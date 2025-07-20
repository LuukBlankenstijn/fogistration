package repository

import (
	"context"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
)

type ContestRepository struct {
	queries *database.Queries
}

func NewContestRepository(queryies *database.Queries) *ContestRepository {
	return &ContestRepository{
		queries: queryies,
	}
}

func (r *ContestRepository) GetNextOrActive(ctx context.Context) (*database.Contest, error) {
	contest, err := r.queries.GetNextOrActiveContest(ctx)
	if len(contest) == 0 {
		return nil, nil
	}

	return &contest[0], err
}

func (r *ContestRepository) GetHashes(ctx context.Context) ([]database.GetContestHashesRow, error) {
	return r.queries.GetContestHashes(ctx)
}

func (r *ContestRepository) Upsert(ctx context.Context, params database.UpsertContestParams) error {
	return r.queries.UpsertContest(ctx, params)
}

func (r *ContestRepository) DeleteAllTeams(ctx context.Context, contestId int32) error {
	return r.queries.DeleteAllContestTeams(ctx, contestId)
}

func (r *ContestRepository) InsertContestTeams(ctx context.Context, contestId int32, teams []database.Team) error {
	params := []database.InsertContestTeamsParams{}
	for _, team := range teams {
		params = append(params, database.InsertContestTeamsParams{
			ContestID: contestId,
			TeamID:    team.ID,
		})
	}
	_, err := r.queries.InsertContestTeams(ctx, params)
	return err
}
