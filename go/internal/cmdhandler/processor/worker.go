package processor

import (
	"context"

	"github.com/LuukBlankenstijn/fogistration/internal/cmdhandler/client/wrapper"
	syncer "github.com/LuukBlankenstijn/fogistration/internal/cmdhandler/sync"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	dbObject "github.com/LuukBlankenstijn/fogistration/internal/shared/database/object"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommandHandler struct {
	client  *wrapper.Client
	sync    *syncer.DomJudgeSyncer
	queries *database.Queries
}

func NewCommandHandler(
	ctx context.Context,
	db *pgxpool.Pool,
	client *wrapper.Client,
) *CommandHandler {
	queries := database.New(db)
	s := syncer.NewSyncer(ctx, client, db)
	return &CommandHandler{
		client:  client,
		sync:    s,
		queries: queries,
	}
}

func (c *CommandHandler) Start(ctx context.Context, dbURL string) {
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		logging.Fatal("failed to connect for notifications", err)
	}
	defer func() {
		err := conn.Close(ctx)
		if err != nil {
			logging.Error("failed to close database connection", err)
		}
	}()

	_, err = conn.Exec(ctx, "LISTEN new_command")
	if err != nil {
		logging.Fatal("failed to listen for commands", err)
	}

	for {
		select {
		case <-ctx.Done():
			return
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

func (c *CommandHandler) tryProcessCommand(ctx context.Context) bool {
	cmdRepo := repository.NewCommandRepository(c.queries)
	cmd, found, err := cmdRepo.TryDequeue(ctx)
	if err != nil || !found {
		return false
	}

	go c.processCommand(ctx, cmd)
	return true
}

func (c *CommandHandler) processCommand(ctx context.Context, cmd dbObject.DatabaseObject) {
	switch typedCmd := cmd.(type) {
	case dbObject.SyncDj:
		c.doSync(c.sync)
	case dbObject.ChangeIp:
		err := c.handleSetIpCommand(ctx, typedCmd)
		if err != nil {
			logging.Error("failed to set ip", err, typedCmd.Id)
		}
	default:
		logging.Error("unknown processor type", nil)
	}
}
