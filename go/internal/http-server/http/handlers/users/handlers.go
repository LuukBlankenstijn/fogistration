package users

import (
	"context"

	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/middleware"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/sse"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/models"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/danielgtaylor/huma/v2"
)

func (h *Handlers) GetCurrentUser(ctx context.Context, request *struct{}) (*getUserResponse, error) {
	user, ok := middleware.User(ctx)
	if !ok {
		return nil, huma.Error404NotFound("user not found")
	}
	return &getUserResponse{Body: models.MapUsers(user)[0]}, nil
}

func (h *Handlers) GetUser(ctx context.Context, request *getUserRequest) (*sse.GetResponse[models.User], error) {
	user, err := h.Q.GetUserByID(ctx, int64(request.ID))
	if err != nil {
		return nil, huma.Error404NotFound("user not found")
	}
	return &sse.GetResponse[models.User]{Body: models.MapUsers(user)[0]}, nil
}

func (h *Handlers) ListUsers(ctx context.Context, request *struct{}) (*listUsersResponse, error) {
	users, err := h.Q.ListUsers(ctx)
	if err != nil {
		users = []database.User{}
	}

	return &listUsersResponse{models.MapUsers(users...)}, nil
}

func (h *Handlers) PutUser(ctx context.Context, request *putUserRequest) (*getUserResponse, error) {
	user, err := h.Q.UpdateUserProfile(ctx, database.UpdateUserProfileParams{
		Username: database.PgTextFromString(&request.Body.Username),
		Email:    database.PgTextFromString(&request.Body.Email),
		Role:     request.Body.Role,
		ID:       int64(request.ID),
	})
	if err != nil {
		logging.Error("failed to update user in db", err)
		return nil, huma.Error500InternalServerError("database error")
	}

	return &getUserResponse{models.MapUsers(user)[0]}, nil
}

func (h *Handlers) patchUser(ctx context.Context, request *patchUserRequest) (*getUserResponse, error) {
	currentUserResponse, err := h.GetUser(ctx, &getUserRequest{ID: request.ID})
	if err != nil {
		return nil, err
	}
	user := models.PatchUser(currentUserResponse.Body, request.Body)

	return h.PutUser(ctx, &putUserRequest{
		ID:   request.ID,
		Body: user,
	})
}
