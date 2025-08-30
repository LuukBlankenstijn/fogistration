package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database/dblisten/internal/common"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Listener struct {
	pool *pgxpool.Pool

	mu  sync.RWMutex
	reg map[string]reflect.Type

	dmx   *common.Demux
	batch int
	every time.Duration

	once sync.Once
}

func New(ctx context.Context, url string, batch int, every time.Duration) (*Listener, error) {
	c, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	l := &Listener{pool: c, reg: make(map[string]reflect.Type), batch: batch, every: every}
	if err := l.ensureInfra(ctx); err != nil {
		c.Close()
		return nil, err
	}
	mixed, err := l.consume(ctx, batch, every)
	if err != nil {
		c.Close()
		return nil, err
	}
	l.dmx = common.NewDemux(mixed)
	return l, nil
}

func (l *Listener) Close(ctx context.Context) { l.pool.Close() }

func (l *Listener) ensureInfra(ctx context.Context) error {
	var err error
	l.once.Do(func() { err = l.EnsureQueueInfra(ctx) })
	return err
}

func (l *Listener) Register(ctx context.Context, table string, typ any) error {
	if err := l.ensureInfra(ctx); err != nil {
		return err
	}
	schema, tbl := common.NormalizeTable(table)
	t := reflect.TypeOf(typ)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("typ must be a struct")
	}
	if err := l.EnsureQueueTrigger(ctx, schema, tbl); err != nil {
		return err
	}
	l.mu.Lock()
	l.reg[schema+"."+tbl] = t
	l.mu.Unlock()
	return nil
}

func (l *Listener) Subscribe(table string, buf int) <-chan common.Mixed {
	return l.dmx.Subscribe(common.ToKey(table), buf)
}

func (l *Listener) consume(ctx context.Context, batch int, catchupEvery time.Duration) (<-chan common.Mixed, error) {
	c, err := l.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	out := make(chan common.Mixed, batch)

	if _, err := c.Conn().Exec(ctx, `LISTEN row_change_queue`); err != nil {
		return nil, err
	}

	drainOnce := func() bool {
		tx, err := l.pool.Begin(ctx)
		if err != nil {
			return false
		}
		rows, err := tx.Query(ctx, `
SELECT id, tbl, op, new_json, old_json
FROM public.row_change_queue
ORDER BY id
FOR UPDATE SKIP LOCKED
LIMIT $1`, batch)
		if err != nil {
			_ = tx.Rollback(ctx)
			return false
		}

		type item struct {
			ID             int64
			Tbl, Op        string
			NewRaw, OldRaw json.RawMessage
		}
		var items []item
		for rows.Next() {
			var it item
			if rows.Scan(&it.ID, &it.Tbl, &it.Op, &it.NewRaw, &it.OldRaw) == nil {
				items = append(items, it)
			}
		}
		rows.Close()
		if len(items) == 0 {
			_ = tx.Rollback(ctx)
			return false
		}

		var ids []int64
		for _, it := range items {
			ids = append(ids, it.ID)
			key := common.ToKey(it.Tbl)
			l.mu.RLock()
			t, ok := l.reg[key]
			l.mu.RUnlock()
			if !ok {
				continue
			}

			var newPtr, oldPtr any
			if len(it.NewRaw) > 0 && string(it.NewRaw) != "null" {
				np := reflect.New(t).Interface()
				if json.Unmarshal(it.NewRaw, np) == nil {
					newPtr = np
				}
			}
			if len(it.OldRaw) > 0 && string(it.OldRaw) != "null" {
				op := reflect.New(t).Interface()
				if json.Unmarshal(it.OldRaw, op) == nil {
					oldPtr = op
				}
			}
			select {
			case out <- common.Mixed{Table: key, Op: it.Op, New: newPtr, Old: oldPtr}:
			case <-ctx.Done():
				_ = tx.Rollback(ctx)
				return true
			}
		}

		if _, err := tx.Exec(ctx, `DELETE FROM public.row_change_queue WHERE id = ANY($1)`, ids); err != nil {
			_ = tx.Rollback(ctx)
			return false
		}
		_ = tx.Commit(ctx)
		return true
	}

	go func() {
		defer close(out)
		defer func() {
			_, _ = c.Conn().Exec(context.Background(), `UNLISTEN row_change_queue`)
			c.Release()
		}()

		_ = drainOnce()
		ticker := time.NewTicker(catchupEvery)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for drainOnce() {
				}
			default:
				_, _ = c.Conn().WaitForNotification(ctx)
				for drainOnce() {
				}
			}
		}
	}()

	return out, nil
}
