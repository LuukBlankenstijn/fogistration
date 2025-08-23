package processor

import (
	"context"
	"time"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database/dblisten"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
)

func (w *Worker) startDBListen(ctx context.Context, dbURL string) error {
	l, err := dblisten.New(ctx, dbURL)
	if err != nil {
		logging.Fatal("failed to created database listener", err)
	}
	defer l.Close(ctx)

	err = l.EnsureQueueInfra(ctx)
	if err != nil {
		logging.Error("failed to ensure infra", err)
		return err
	}

	err = l.RegisterQueue("teams", database.Team{})
	if err != nil {
		logging.Error("failed to register Queue", err)
		return err
	}

	mixed, err := l.ConsumeQueueWithNotify(ctx, 5, 30*time.Second)
	if err != nil {
		logging.Error("failed to get Queue channel", err)
		return err
	}

	for team_update := range dblisten.View[database.Team]("teams", mixed) {
		err := w.handleSetIpCommand(ctx, team_update)
		if err != nil {
			logging.Error("failed to set ip", err)
		}
	}

	return nil
}
