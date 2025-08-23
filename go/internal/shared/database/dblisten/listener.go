package dblisten

import (
	"context"
	"reflect"

	"github.com/jackc/pgx/v5"
)

type Listener struct {
	conn *pgx.Conn
	reg  map[string]reflect.Type
}

func New(ctx context.Context, url string) (*Listener, error) {
	c, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, err
	}
	return &Listener{conn: c, reg: make(map[string]reflect.Type)}, nil
}

func (l *Listener) Close(ctx context.Context) error {
	return l.conn.Close(ctx)
}

func (l *Listener) RegisterQueue(table string, typ any) error {
	schema, tbl := normalizeTable(table)
	t := reflect.TypeOf(typ)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return ErrNotStruct
	}
	if err := l.EnsureQueueTrigger(context.Background(), schema, tbl); err != nil {
		return err
	}
	l.reg[schema+"."+tbl] = t
	return nil
}

func (l *Listener) RegisterNotify(ctx context.Context, table string, typ any) error {
	schema, tbl := normalizeTable(table)
	t := reflect.TypeOf(typ)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return ErrNotStruct
	}
	err := l.EnsureNotifyTrigger(ctx, schema, tbl)
	if err != nil {
		return err
	}
	l.reg[schema+"."+tbl] = t
	return nil
}
