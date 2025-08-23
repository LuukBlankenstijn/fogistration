package processor

import (
	"context"

	"github.com/LuukBlankenstijn/fogistration/internal/cmdhandler/client/wrapper"
	syncer "github.com/LuukBlankenstijn/fogistration/internal/cmdhandler/sync"
	"golang.org/x/sync/errgroup"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Worker struct {
	client  *wrapper.Client
	sync    *syncer.DomJudgeSyncer
	queries *database.Queries
}

func NewWorker(
	ctx context.Context,
	db *pgxpool.Pool,
	client *wrapper.Client,
) *Worker {
	queries := database.New(db)
	s := syncer.NewSyncer(ctx, client, db)
	return &Worker{
		client:  client,
		sync:    s,
		queries: queries,
	}
}

func (c *Worker) Start(ctx context.Context, dbURL string) {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return c.startCommandHandler(ctx, dbURL)
	})

	g.Go(func() error {
		return c.startDBListen(ctx, dbURL)
	})

	if err := g.Wait(); err != nil && err != context.Canceled {
		logging.Fatal("", err)
	}
}
