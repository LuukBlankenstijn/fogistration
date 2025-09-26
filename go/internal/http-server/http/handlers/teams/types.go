package teams

import "github.com/LuukBlankenstijn/fogistration/internal/http-server/models"

type teamResponse struct {
	Body models.Team
}

type getTeamRequest struct {
	ID string `path:"id" doc:"Team external ID"`
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

type getPrintInfoRequest struct {
	Ip string `query:"ip"`
}

type getPrintInfoResponse struct {
	Body struct {
		Name string `json:"teamname"`
	}
}
