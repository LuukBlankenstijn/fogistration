package users

import (
	"context"

	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/middleware"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/sse"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/models"
	"github.com/danielgtaylor/huma/v2"
)

func (h *Handlers) GetCurrentUser(ctx context.Context, request *struct{}) (*sse.GetResponse[models.User], error) {
	user, ok := middleware.User(ctx)
	if !ok {
		return nil, huma.Error404NotFound("user not found")
	}
	return &sse.GetResponse[models.User]{Body: models.MapUser(user)}, nil
}
