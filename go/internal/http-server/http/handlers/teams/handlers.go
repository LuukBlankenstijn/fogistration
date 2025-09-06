package teams

import (
	"context"
	"errors"

	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/sse"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/models"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/service"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/danielgtaylor/huma/v2"
)

func (h *Handlers) getSingleTeam(ctx context.Context, req *getTeamRequest) (*sse.GetResponse[models.Team], error) {
	team, err := h.Q.GetTeamByExternalId(ctx, req.ID)
	if err != nil {
		return nil, huma.Error400BadRequest("team not found")
	}

	return &sse.GetResponse[models.Team]{Body: models.MapTeam(team)[0]}, nil
}

func (h *Handlers) listTeams(ctx context.Context, req *struct{}) (*listTeamsResponse, error) {
	teams, err := h.Q.GetAllTeams(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("could not get teams")
	}

	return &listTeamsResponse{
		Body: models.MapTeam(teams...),
	}, nil
}

func (h *Handlers) setClient(ctx context.Context, request *setClientRequest) (*teamResponse, error) {
	// Load team
	team, err := h.Q.GetTeamByExternalId(ctx, request.ID)
	if err != nil {
		return nil, huma.Error400BadRequest("team not found")
	}

	// Load client if provided
	var client *database.Client
	if request.Body.ClientID != nil {
		c, err := h.Q.GetClientById(ctx, *request.Body.ClientID)
		if err != nil {
			return nil, huma.Error400BadRequest("client not found")
		}
		client = &c
	}

	if err := h.S.Team.SetClient(ctx, team, client); err != nil {
		switch {
		case errors.Is(err, service.ErrClientIPUsedByOtherTeam):
			return nil, huma.Error400BadRequest("client is already associated to another team")
		default:
			logging.Error("SetTeamClient failed", err, "team", team.ExternalID, "clientId", request.Body.ClientID)
			return nil, huma.Error500InternalServerError("internal error")
		}
	}

	// Return updated team (re-fetch for current state)
	updated, err := h.Q.GetTeamByExternalId(ctx, team.ExternalID)
	if err != nil {
		logging.Error("refetch updated team failed", err, "team", team.ExternalID)
		return nil, huma.Error500InternalServerError("internal error")
	}

	return &teamResponse{Body: models.MapTeam(updated)[0]}, nil
}
