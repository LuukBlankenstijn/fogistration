package users

import (
	"net/http"

	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/container"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/middleware"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/sse"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database/models"
	"github.com/danielgtaylor/huma/v2"
)

type Handlers struct {
	*container.Container
}

func NewHandlers(container *container.Container) *Handlers {
	return &Handlers{container}
}

func (h *Handlers) Register(
	api huma.API,
	middlewareFactory *middleware.MiddlewareFactory,
	prefixes ...string,
) {
	userApi := huma.NewGroup(api, prefixes...)

	userApi.UseMiddleware(middlewareFactory.Auth(userApi)...)

	huma.Register(userApi, huma.Operation{
		OperationID: "getCurrentUser",
		Method:      http.MethodGet,
		Path:        "/me",
		Summary:     "Get the current logged in user",
		Tags:        []string{"users"},
	}, h.GetCurrentUser)

	adminGroup := huma.NewGroup(userApi)
	adminGroup.UseMiddleware(middlewareFactory.RequireRoles(adminGroup, models.Admin))

	sse.Register(h.SSE, userApi, huma.Operation{
		OperationID: "getUser",
		Method:      http.MethodGet,
		Path:        "/{id}",
		Summary:     "Get a user by id",
		Tags:        []string{"users"},
	}, h.GetUser)

	huma.Register(userApi, huma.Operation{
		OperationID: "listUsers",
		Method:      http.MethodGet,
		Path:        "/",
		Summary:     "Get all users",
		Tags:        []string{"users"},
	}, h.ListUsers)

	huma.Register(userApi, huma.Operation{
		OperationID: "putUser",
		Method:      http.MethodPut,
		Path:        "/{id}",
		Summary:     "Create or replace a user by id",
		Tags:        []string{"users"},
	}, h.PutUser)

	huma.Register(userApi, huma.Operation{
		OperationID: "patchUser",
		Method:      http.MethodPatch,
		Path:        "/{id}",
		Summary:     "Patch a user",
		Tags:        []string{"users"},
	}, h.patchUser)
}
