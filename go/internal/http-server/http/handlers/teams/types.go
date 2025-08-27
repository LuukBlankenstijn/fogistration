package teams

import "github.com/LuukBlankenstijn/fogistration/internal/http-server/models"

type teamResponse struct {
	Body models.Team
}

type listTeamsResponse struct {
	Body []models.Team `nullable:"false"`
}

type setClientRequest struct {
	Body struct {
		ClientID *int32 `json:"clientId,omitempty"`
	}
	ID string `path:"id" doc:"Team external ID"`
}
