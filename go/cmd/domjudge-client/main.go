package main

import (
	"context"

	syncer "github.com/LuukBlankenstijn/fogistration/internal/domjudge/sync"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var cfg config.DomJudgeConfig
	err := config.Load(&cfg, ".env-dj")
	if err != nil {
		logging.Fatal("Failed to load config: %v", err)
	}

	logging.SetupLogger(cfg.LogLevel, cfg.AppEnv)

	ctx := context.Background()
	url := database.GetUrl(&cfg.DB)
	dbpool, err := pgxpool.New(ctx, url)
	if err != nil {
		logging.Fatal("unable to create dbpool: %w", err)
	}
	defer dbpool.Close()

	syncer, err := syncer.NewSyncer(ctx, cfg, dbpool)
	if err != nil {
		logging.Fatal("unable to create syncer: %w", err)
	}

	err = syncer.Sync()
	if err != nil {
		logging.Fatal("unable to sync: %w", err)
	}

}
