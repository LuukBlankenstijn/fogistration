package dblisten

import (
	"context"
	"fmt"
	"time"

	inotify "github.com/LuukBlankenstijn/fogistration/internal/shared/database/dblisten/internal/notify"
	iqueue "github.com/LuukBlankenstijn/fogistration/internal/shared/database/dblisten/internal/queue"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database/dblisten/internal/common"
)

// Re-export types so users only import dblisten.
type Mixed = common.Mixed
type Notification[T any] = common.Notification[T]

// Listener is implemented by both modes but kept minimal.
type Listener interface {
	Close(ctx context.Context)
	// hidden core methods accessed via type assertion internally
}

type listenerCore interface {
	Register(ctx context.Context, table string, typ any) error
	Subscribe(table string, buf int) <-chan common.Mixed
}

// Constructors auto-ensure infra and start streaming.
func NewNotify(ctx context.Context, url string) (Listener, error) {
	l, err := inotify.New(ctx, url)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func NewQueue(ctx context.Context, url string, batch int, every time.Duration) (Listener, error) {
	l, err := iqueue.New(ctx, url, batch, every)
	if err != nil {
		return nil, err
	}
	return l, nil
}

// RegisterTyped ensures trigger and returns typed channel.
func RegisterTyped[T any](ctx context.Context, l Listener, table string, buf int) (<-chan Notification[T], error) {
	core, ok := l.(listenerCore)
	if !ok {
		return nil, fmt.Errorf("unsupported listener")
	}
	if err := core.Register(ctx, table, new(T)); err != nil {
		return nil, err
	}
	raw := core.Subscribe(table, buf)
	return common.View[T](table, raw), nil
}
