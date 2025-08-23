package processor

import (
	"context"
	"fmt"

	dbObject "github.com/LuukBlankenstijn/fogistration/internal/shared/database/command"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/repository"
	"github.com/jackc/pgx/v5"
)

func (c *Worker) startCommandHandler(ctx context.Context, dbURL string) error {
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect for notifications: %w", err)
	}
	defer func() {
		err := conn.Close(ctx)
		if err != nil {
			logging.Error("failed to close database connection", err)
		}
	}()

	_, err = conn.Exec(ctx, "LISTEN new_command")
	if err != nil {
		return fmt.Errorf("failed to listen for commands: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if c.tryProcessCommand(ctx) {
				continue // Check for more commands
			}
			_, err := conn.WaitForNotification(ctx)
			if err != nil {
				logging.Error("error while waiting for notifications", err)
			}
		}
	}
}

func (c *Worker) tryProcessCommand(ctx context.Context) bool {
	cmdRepo := repository.NewCommandRepository(c.queries)
	cmd, found, err := cmdRepo.TryDequeue(ctx)
	if err != nil || !found {
		return false
	}

	go c.processCommand(cmd)
	return true
}

func (c *Worker) processCommand(cmd dbObject.Command) {
	switch cmd.(type) {
	case dbObject.SyncDj:
		c.doSync(c.sync)
	default:
		logging.Error("unknown processor type", nil)
	}
}
