package processor

import (
	"context"

	"github.com/LuukBlankenstijn/fogistration/internal/cmdhandler/client/wrapper"
	syncer "github.com/LuukBlankenstijn/fogistration/internal/cmdhandler/sync"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/command"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommandHandler struct {
	cmdRepo  *repository.CommandRepository
	teamRepo *repository.TeamRepository
	client   *wrapper.Client
	sync     *syncer.DomJudgeSyncer
}

func NewCommandHandler(
	ctx context.Context,
	db *pgxpool.Pool,
	client *wrapper.Client,
) *CommandHandler {
	queries := database.New(db)
	s := syncer.NewSyncer(ctx, client, db)
	return &CommandHandler{
		cmdRepo:  repository.NewCommandRepository(queries),
		teamRepo: repository.NewTeamRepository(queries),
		client:   client,
		sync:     s,
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
	cmd, found, err := c.cmdRepo.TryDequeue(ctx)
	if err != nil || !found {
		return false
	}

	go c.processCommand(ctx, cmd)
	return true
}

func (c *CommandHandler) processCommand(ctx context.Context, cmd command.Command) {
	switch typedCmd := cmd.(type) {
	case command.SyncDjCommand:
		c.doSync(c.sync)
	case command.SetIpCommand:
		err := c.handleSetIpCommand(ctx, typedCmd)
		if err != nil {
			logging.Error("failed to set ip", err, typedCmd.Id)
		}
	default:
		logging.Error("unknown processor type", nil)
	}
}
