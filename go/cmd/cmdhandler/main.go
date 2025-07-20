package main

import (
	"context"
	"time"

	"github.com/LuukBlankenstijn/fogistration/internal/cmdhandler/client/wrapper"

	"github.com/LuukBlankenstijn/fogistration/internal/cmdhandler/processor"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/command"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var cfg config.DomJudgeConfig
	err := config.Load(&cfg, ".env-cmdhandler")
	if err != nil {
		logging.Fatal("Failed to load config", err)
	}

	logging.SetupLogger(cfg.LogLevel, cfg.AppEnv)

	ctx := context.Background()
	url := database.GetUrl(&cfg.DB)
	dbpool, err := pgxpool.New(ctx, url)
	if err != nil {
		logging.Fatal("unable to create dbpool", err)
	}
	defer dbpool.Close()

	djClient, err := wrapper.NewClient(ctx, cfg)
	if err != nil {
		logging.Fatal("failed to create cmdhandler client", err)
	}
	commandHandler := processor.NewCommandHandler(ctx, dbpool, djClient)

	// TODO: move this somewhere else
	queries := database.New(dbpool)
	interval, err := time.ParseDuration(cfg.SyncInterval)
	if err != nil {
		logging.Fatal("unable to parse interval", err)
	}
	go startScheduledSync(ctx, interval, queries)

	commandHandler.Start(ctx, url)
}

func startScheduledSync(
	ctx context.Context,
	interval time.Duration,
	queries *database.Queries,
) {
	cmdRepo := repository.NewCommandRepository(queries)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := cmdRepo.Enqueue(ctx, command.SyncDjCommand{})
			if err != nil {
				logging.Error("failed to enqueue SyncDjCommand", err)
			}
		}
	}
}
