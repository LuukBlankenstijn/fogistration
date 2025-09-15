package main

import (
	"context"

	httpServer "github.com/LuukBlankenstijn/fogistration/internal/http-server/http"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/seeder"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var cfg config.HttpConfig
	ctx := context.Background()
	err := config.Load(&cfg, ".env-http")
	if err != nil {
		logging.Fatal("Failed to load config: %v", err)
	}

	logging.SetupLogger(cfg.LogLevel, cfg.AppEnv)

	url := database.GetUrl(&cfg.DB)
	dbpool, err := pgxpool.New(ctx, url)
	if err != nil {
		logging.Fatal("unable to create dbpool: %w", err)
	}
	defer dbpool.Close()
	queries := database.New(dbpool)

	if cfg.AppEnv == "development" {
		sdr := seeder.New(dbpool, queries)
		err = sdr.SeedDefaultUser(ctx)
		if err != nil {
			logging.Error("seeder failed", err)
		}
	}

	httpServer.NewServer(&cfg, dbpool).Run(ctx)
}
