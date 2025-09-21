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
		Summary:       "OIDC login",
	}, h.handleLogin)

	huma.Register(oidcApi, huma.Operation{
		OperationID:   "oidcCallback",
		Method:        http.MethodGet,
		Path:          "/callback",
		Tags:          []string{"auth"},
		DefaultStatus: 302,
		Summary:       "OIDC callback",
	}, h.handleCallback)

	huma.Register(oidcApi, huma.Operation{
		OperationID: "oidcEnabled",
		Method:      http.MethodGet,
		Path:        "/enabled",
		Tags:        []string{"auth"},
		Summary:     "OIDC enabled check",
		Description: "Returns 204 when oidc is enabled, 404 otherwise",
	}, h.handleIsEnabled)

}
