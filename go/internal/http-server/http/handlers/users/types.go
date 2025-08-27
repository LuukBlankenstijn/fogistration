package users

import "github.com/LuukBlankenstijn/fogistration/internal/http-server/models"

type GetCurrentUserResponse struct {
	Body models.User
}
