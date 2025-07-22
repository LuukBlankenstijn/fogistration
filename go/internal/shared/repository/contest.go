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
