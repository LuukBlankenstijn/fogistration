package main

import (
	"context"

	"github.com/LuukBlankenstijn/fogistration/internal/grpc/pubsub"
	"github.com/LuukBlankenstijn/fogistration/internal/grpc/server"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	var cfg config.GrpcConfig
	err := config.Load(&cfg, ".env-grpc")
	if err != nil {
		logging.Fatal("Failed to load config", err)
	}

	url := database.GetUrl(&cfg.DB)
	dbpool, err := pgxpool.New(ctx, url)
	if err != nil {
		logging.Fatal("unable to create dbpool", err)
	}
	defer dbpool.Close()

	queries := database.New(dbpool)
	pubsub := pubsub.NewManager()
	server := server.NewServer(queries, pubsub)

	if err := server.Start(cfg.Port); err != nil {
		logging.Error("gRPC server failed", err)
	}
}
