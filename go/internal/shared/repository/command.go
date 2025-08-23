package repository

import (
	"context"
	"encoding/json"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	dbObject "github.com/LuukBlankenstijn/fogistration/internal/shared/database/command"
	"github.com/jackc/pgx/v5"
)

type CommandRepository struct {
	queries *database.Queries
}

func NewCommandRepository(queries *database.Queries) *CommandRepository {
	return &CommandRepository{
		queries: queries,
	}
}

func (r *CommandRepository) Enqueue(ctx context.Context, cmd dbObject.Command) error {
	payload, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	return r.queries.EnqueueCommand(ctx, database.EnqueueCommandParams{
		CommandType: string(cmd.Type()),
		Payload:     payload,
	})
}

func (r *CommandRepository) Dequeue(ctx context.Context) (dbObject.Command, error) {
	row, err := r.queries.DequeueCommand(ctx)
	if err != nil {
		return nil, err
	}

	return dbObject.ParseDatabaseObject(row.CommandType, row.Payload)
}

func (r *CommandRepository) TryDequeue(ctx context.Context) (dbObject.Command, bool, error) {
	cmd, err := r.Dequeue(ctx)
	if err == pgx.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return cmd, true, nil
}
