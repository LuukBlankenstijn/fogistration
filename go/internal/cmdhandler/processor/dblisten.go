package processor

import (
	"context"
	"fmt"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database/dblisten"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
)

func (w *Worker) startDBListen(ctx context.Context, dbURL string) error {
	l, err := dblisten.NewNotify(ctx, dbURL)
	if err != nil {
		return fmt.Errorf("failed to create database listener: %w", err)
	}
	defer l.Close(ctx)

	teams, err := dblisten.RegisterTyped[database.Team](ctx, l, "teams", 32)
	if err != nil {
		return fmt.Errorf("failed to register team listener: %w", err)
	}

	for t := range teams {
		err := w.handleSetIpCommand(ctx, t)
		if err != nil {
			logging.Error("failed to set ip", err)
		}
	}

	return nil
}
