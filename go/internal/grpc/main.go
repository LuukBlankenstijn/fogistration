package grpc

import (
	"context"
	"fmt"

	"github.com/LuukBlankenstijn/fogistration/internal/grpc/notifications"
	"github.com/LuukBlankenstijn/fogistration/internal/grpc/pubsub"
	"github.com/LuukBlankenstijn/fogistration/internal/grpc/server"
	"github.com/LuukBlankenstijn/fogistration/internal/grpc/service"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, cfg config.GrpcConfig) error {
	dbUrl := database.GetUrl(&cfg.DB)
	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return fmt.Errorf("unable to create dbpool: %w", err)
	}

	queries := database.New(pool)
	pubsub := pubsub.NewManager()
	server := server.NewServer(queries, pubsub, cfg)
	serviceContainer := service.New(queries, pubsub, cfg)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return notifications.New(serviceContainer, queries).Run(ctx, dbUrl)
	})

	g.Go(func() error {
		return server.Start()
	})

	return g.Wait()
}
