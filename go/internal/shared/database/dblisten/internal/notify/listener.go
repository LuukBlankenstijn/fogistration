package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database/dblisten/internal/common"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Listener struct {
	pool *pgxpool.Pool

	mu  sync.RWMutex
	reg map[string]reflect.Type

	dmx  *common.Demux
	once sync.Once
}

func New(ctx context.Context, url string) (*Listener, error) {
	c, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	l := &Listener{pool: c, reg: make(map[string]reflect.Type)}
	if err := l.ensureInfra(ctx); err != nil {
		c.Close()
		return nil, err
	}
	mixed, err := l.listen(ctx)
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
	l.once.Do(func() { err = l.EnsureNotifyInfra(ctx) })
	return err
}

// --- core API used by facade ---

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
	if err := l.EnsureNotifyTrigger(ctx, schema, tbl); err != nil {
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

// --- listen -> Mixed ---

func (l *Listener) listen(ctx context.Context) (<-chan common.Mixed, error) {
	c, err := l.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := c.Conn().Exec(ctx, `LISTEN row_changes`); err != nil {
		c.Release()
		return nil, err
	}
	out := make(chan common.Mixed, 256)
	go func() {
		defer close(out)
		for {
			ntf, err := c.Conn().WaitForNotification(ctx)
			if err != nil {
				return
			}
			var p struct {
				Table string          `json:"table"`
				Op    string          `json:"op"`
				New   json.RawMessage `json:"new"`
				Old   json.RawMessage `json:"old"`
			}
			if json.Unmarshal([]byte(ntf.Payload), &p) != nil {
				continue
			}
			key := common.ToKey(p.Table)

			l.mu.RLock()
			t, ok := l.reg[key]
			l.mu.RUnlock()
			if !ok {
				continue
			}

			var newPtr, oldPtr any
			if len(p.New) > 0 && string(p.New) != "null" {
				v := reflect.New(t).Interface()
				if json.Unmarshal(p.New, v) == nil {
					newPtr = v
				}
			}
			if len(p.Old) > 0 && string(p.Old) != "null" {
				v := reflect.New(t).Interface()
				if json.Unmarshal(p.Old, v) == nil {
					oldPtr = v
				}
			}

			select {
			case out <- common.Mixed{Table: key, Op: p.Op, New: newPtr, Old: oldPtr}:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, nil
}
