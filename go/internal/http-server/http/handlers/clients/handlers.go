package clients

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

func (h *Handlers) getSingleClient(ctx context.Context, req *getClientRequest) (*sse.GetResponse[models.Client], error) {
	client, err := h.Q.GetClientById(ctx, req.ID)
	if err != nil {
		return nil, huma.Error400BadRequest("team not found")
	}

	return &sse.GetResponse[models.Client]{Body: models.MapClient(client)[0]}, nil
}

func (h *Handlers) listClients(ctx context.Context, request *struct{}) (*listClientsResponse, error) {
	clients, err := h.Q.GetAllClients(ctx)
	if err != nil {
		logging.Error("failed to get all clients", err)
		return nil, huma.Error500InternalServerError("could not get clients")
	}

	return &listClientsResponse{models.MapClient(clients...)}, nil
}

func (h *Handlers) setTeam(ctx context.Context, request *setTeamRequest) (*setTeamResponse, error) {
	// Load client
	cl, err := h.Q.GetClientById(ctx, request.ID)
	if err != nil {
		return nil, huma.Error400BadRequest("client not found")
	}

	// Load team if provided
	var team *database.Team
	if request.Body.TeamID != nil {
		t, err := h.Q.GetTeamByExternalId(ctx, *request.Body.TeamID)
		if err != nil {
			return nil, huma.Error400BadRequest("team not found")
		}
		team = &t
	}

	// Delegate to service (owns TX)
	if err := h.S.Client.AssignClientToTeam(ctx, cl, team); err != nil {
		switch {
		case errors.Is(err, service.ErrTeamAlreadyAssigned):
			return nil, huma.Error400BadRequest("new team already has a client assigned")
		default:
			logging.Error("AssignClientToTeam failed", err, "clientID", cl.ID, "teamID", request.Body.TeamID)
			return nil, huma.Error500InternalServerError("internal error")
		}
	}

	return &setTeamResponse{Body: models.MapClient(cl)[0]}, nil
}
