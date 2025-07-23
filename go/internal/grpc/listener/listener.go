package listener

import (
	"context"
	"fmt"

	"github.com/LuukBlankenstijn/fogistration/internal/grpc/pubsub"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseListener struct {
	pubsub  *pubsub.Manager
	db      *pgxpool.Pool
	queries *database.Queries
}

func NewDatabaseListener(pubsub *pubsub.Manager, db *pgxpool.Pool) *DatabaseListener {
	return &DatabaseListener{
		pubsub:  pubsub,
		db:      db,
		queries: database.New(db),
	}
}

func (d *DatabaseListener) Run(ctx context.Context) error {

	conn, err := d.db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	if err := d.setupConnections(ctx, conn); err != nil {
		return fmt.Errorf("failed to register listeners: %w", err)
	}

	d.listen(ctx, conn)

	return nil
}
