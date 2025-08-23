package dblisten

import (
	"context"
	"encoding/json"
	"reflect"
)

func (l *Listener) ListenNotify(ctx context.Context) (<-chan Mixed, error) {
	if _, err := l.conn.Exec(ctx, `LISTEN row_changes`); err != nil {
		return nil, err
	}
	out := make(chan Mixed, 256)

	go func() {
		defer close(out)
		for {
			ntf, err := l.conn.WaitForNotification(ctx)
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

			t, ok := l.reg[toKey(p.Table)]
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
			case out <- Mixed{Table: toKey(p.Table), Op: p.Op, New: newPtr, Old: oldPtr}:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out, nil
}
