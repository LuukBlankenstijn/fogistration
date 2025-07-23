package listener

import (
	"context"
	"fmt"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (d *DatabaseListener) setupConnections(ctx context.Context, db *pgxpool.Conn) error {
	for t := range handlers {
		_, err := db.Exec(ctx, fmt.Sprintf("LISTEN %s", t))
		if err != nil {
			logging.Error(fmt.Sprintf("failed to listen to register listener for type %s", t), err)
		}
		logging.Info("setting up listener for type %s", t)
	}
	return nil
}

func (d *DatabaseListener) listen(ctx context.Context, db *pgxpool.Conn) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			notification, err := db.Conn().WaitForNotification(ctx)
			if err != nil {
				logging.Error("error while waiting for notifications", err)
			}
			err = d.dispatchDatabaseObject(ctx, notification.Channel, []byte(notification.Payload))
			if err != nil {
				logging.Error("failed to dispatch object", err)
			}
		}
	}
}
