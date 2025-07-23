package listener

import (
	"context"
	"fmt"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database/object"
)

type DatabaseObjectHandler func(*DatabaseListener, context.Context, object.DatabaseObject) error

var handlers = map[object.DatabaseObjectType]DatabaseObjectHandler{
	object.IpChangeType: (*DatabaseListener).handleIpChange,
}

func (d *DatabaseListener) dispatchDatabaseObject(ctx context.Context, typ string, payload []byte) error {
	obj, err := object.ParseDatabaseObject(typ, payload)
	if err != nil {
		return fmt.Errorf("failed to parse object: %w", err)
	}

	handler, ok := handlers[object.DatabaseObjectType(typ)]
	if !ok {
		return fmt.Errorf("no handler registered for type %q", typ)
	}

	return handler(d, ctx, obj)
}
