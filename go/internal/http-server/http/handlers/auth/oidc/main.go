package oidc

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
	oidcApi := huma.NewGroup(api, prefixes...)

	huma.Register(oidcApi, huma.Operation{
		OperationID:   "oidcLogin",
		Method:        http.MethodGet,
		Path:          "/login",
		Tags:          []string{"auth"},
		DefaultStatus: 302,
	}, h.handleLogin)

	huma.Register(oidcApi, huma.Operation{
		OperationID:   "oidcCallback",
		Method:        http.MethodGet,
		Path:          "/callback",
		Tags:          []string{"auth"},
		DefaultStatus: 302,
	}, h.handleCallback)

}
