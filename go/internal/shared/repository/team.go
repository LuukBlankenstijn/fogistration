package repository

import (
	"context"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
)

type TeamRepository struct {
	queries *database.Queries
}

func NewTeamRepository(queryies *database.Queries) *TeamRepository {
	return &TeamRepository{
		queries: queryies,
	}
}

func (r *TeamRepository) GetHashes(ctx context.Context) ([]database.GetContestHashesRow, error) {
	return r.queries.GetContestHashes(ctx)
}

func (r *TeamRepository) Upsert(ctx context.Context, params database.UpsertTeamParams) error {
	return r.queries.UpsertTeam(ctx, params)
}
