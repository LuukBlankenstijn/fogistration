package contests

import (
	"context"

	"github.com/LuukBlankenstijn/fogistration/internal/http-server/models"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
)

func (h *Handlers) getAllContests(ctx context.Context, request *struct{}) (*listContestsResponse, error) {
	contests, err := h.Q.ListContests(ctx)
	if err != nil {
		logging.Error("failed to get all contests", err)
		return nil, huma.Error500InternalServerError("could not get contests")
	}

	return &listContestsResponse{models.MapContest(contests...)}, nil
}

func (h *Handlers) getActiveContest(ctx context.Context, request *struct{}) (*contestResponse, error) {
	contest, err := h.Q.GetNextOrActiveContest(ctx)
	if err == pgx.ErrNoRows {
		return nil, huma.Error404NotFound("")
	}
	if err != nil {
		logging.Error("failed to get active contest", err)
		return nil, huma.Error500InternalServerError("could not get active contest")
	}

	return &contestResponse{models.MapContest(contest)[0]}, nil
}
