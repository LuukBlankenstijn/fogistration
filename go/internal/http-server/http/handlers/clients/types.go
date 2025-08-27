package clients

import "github.com/LuukBlankenstijn/fogistration/internal/http-server/models"

type listClientsResponse struct {
	Body []models.Client
}

type setTeamRequest struct {
	ID   int32 `path:"id" doc:"Client ID"`
	Body struct {
		TeamID *string `json:"teamId,omitempty"`
	}
}

type setTeamResponse struct {
	Body models.Client
}
