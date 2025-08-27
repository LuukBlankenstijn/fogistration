package users

import (
	"context"

	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/middleware"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/models"
	"github.com/danielgtaylor/huma/v2"
)

func (h *Handlers) GetCurrentUser(ctx context.Context, request *struct{}) (*GetCurrentUserResponse, error) {
	user, ok := middleware.User(ctx)
	if !ok {
		return nil, huma.Error404NotFound("user not found")
	}
	return &GetCurrentUserResponse{models.MapUser(user)}, nil
}
