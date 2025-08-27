package contests

import "github.com/LuukBlankenstijn/fogistration/internal/http-server/models"

type contestResponse struct {
	Body models.Contest
}

type listContestsResponse struct {
	Body []models.Contest
}
