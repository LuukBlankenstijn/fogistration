package users

import (
	"net/http"

	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/container"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/middleware"
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
		Summary:     "Get Current user",
		Tags:        []string{"users"},
	}, h.GetCurrentUser)
}
