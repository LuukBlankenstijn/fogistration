package dblisten

import (
	"context"
	"encoding/json"
	"reflect"
	"time"
)

func (l *Listener) ConsumeQueueWithNotify(ctx context.Context, batch int, catchupEvery time.Duration) (<-chan Mixed, error) {
	out := make(chan Mixed, batch)

	if _, err := l.conn.Exec(ctx, `LISTEN row_change_queue`); err != nil {
		return nil, err
	}

	drainOnce := func() bool {
		tx, err := l.conn.Begin(ctx)
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
			t, ok := l.reg[toKey(it.Tbl)]
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
			case out <- Mixed{Table: toKey(it.Tbl), Op: it.Op, New: newPtr, Old: oldPtr}:
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
				_, _ = l.conn.WaitForNotification(ctx)
				for drainOnce() {
				}
			}
		}
	}()

	return out, nil
}
