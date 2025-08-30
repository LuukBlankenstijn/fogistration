package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/LuukBlankenstijn/fogistration/internal/grpc/eventhandler"
	"github.com/LuukBlankenstijn/fogistration/internal/grpc/pubsub"
	"github.com/LuukBlankenstijn/fogistration/internal/grpc/server"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database/dblisten"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Load config
	var cfg config.GrpcConfig
	if err := config.Load(&cfg, ".env-grpc"); err != nil {
		logging.Fatal("failed to load config", err)
	}

	// Init DB pool
	url := database.GetUrl(&cfg.DB)
	dbpool, err := pgxpool.New(ctx, url)
	if err != nil {
		logging.Fatal("unable to create dbpool", err)
	}
	defer dbpool.Close()

	// Dependencies
	queries := database.New(dbpool)
	pubsub := pubsub.NewManager()
	eventHandler := eventhandler.New(pubsub)
	srv := server.NewServer(queries, pubsub)

	// Concurrent management
	g, ctx := errgroup.WithContext(ctx)

	// g.Go(func() error {
	// 	for {
	// 		select {
	// 		case <-ctx.Done():
	// 			return ctx.Err()
	// 		default:
	// 		}
	//
	// 		logging.Info("starting database listener...")
	// 		err := dbListener.Run(ctx)
	// 		if err != nil && ctx.Err() == nil {
	// 			logging.Error("database listener crashed, retrying in 5s", err)
	// 			time.Sleep(5 * time.Second)
	// 			continue
	// 		}
	// 		return err
	// 	}
	// })

	g.Go(func() error {
		l, err := dblisten.New(ctx, url)
		if err != nil {
			logging.Fatal("failed to created database listener", err)
		}
		defer l.Close(ctx)

		err = l.EnsureNotifyInfra(ctx)
		if err != nil {
			logging.Error("failed to ensure infra", err)
			return err
		}

		err = l.RegisterNotify(ctx, "teams", database.Team{})
		if err != nil {
			logging.Error("failed to register notify", err)
			return err
		}

		mixed, err := l.ListenNotify(ctx)
		if err != nil {
			logging.Error("failed to get notification channel", err)
			return err
		}

		for team_update := range dblisten.View[database.Team]("teams", mixed) {
			eventHandler.HandleIpChange(team_update)
		}

		return nil
	})

	// Start gRPC server
	g.Go(func() error {
		logging.Info("starting gRPC server...")
		return srv.Start(cfg.Port)
	})

	// Wait for termination or error
	if err := g.Wait(); err != nil && err != context.Canceled {
		logging.Error("fatal shutdown", err)
		os.Exit(1)
	}
	logging.Info("shutdown complete")
}
