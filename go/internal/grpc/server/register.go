package server

import (
	"context"
	"fmt"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func register(ctx context.Context, queries *database.Queries, client database.Client) error {
	sync := func() error {
		return queries.SetPendingSync(ctx, database.SetPendingSyncParams{
			ID:          client.ID,
			PendingSync: true,
		})
	}

	_, err := queries.GetTeamByIp(ctx, pgtype.Text{String: client.Ip, Valid: true})

	// err nill, already registered
	if err == nil {
		return sync()
	}

	if err != pgx.ErrNoRows {
		return fmt.Errorf("error when getting team from database: %w", err)
	}

	contest, err := queries.GetNextOrActiveContest(ctx)
	if err != nil {
		return fmt.Errorf("error when getting active contest form database")
	}
	_, err = queries.ClaimTeam(ctx, database.ClaimTeamParams{
		Ip:        database.PgTextFromString(&client.Ip),
		ContestID: contest.ID,
	})
	if err == pgx.ErrNoRows {
		// no team available, do nothing
		return sync()
	}
	if err != nil {
		return fmt.Errorf("error when claiming team: %w", err)
	}

	return sync()
}
