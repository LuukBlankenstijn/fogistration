package grpc

import (
	"context"
	"fmt"

	"github.com/LuukBlankenstijn/fogistration/internal/grpc/notifications"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Run(ctx context.Context, cfg config.GrpcConfig) error {

	dbUrl := database.GetUrl(&cfg.DB)
	_, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return fmt.Errorf("unable to create dbpool: %w", err)
	}

	_ = notifications.Run(ctx, dbUrl)

	return nil
}
