package auth

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
	authAPI := huma.NewGroup(api, prefixes...)

	huma.Register(authAPI, huma.Operation{
		OperationID: "login",
		Method:      http.MethodPost,
		Path:        "/login",
		Summary:     "Login",
		Tags:        []string{"auth"},
	}, h.login)

	huma.Register(authAPI, huma.Operation{
		OperationID:   "logout",
		Method:        http.MethodPost,
		Path:          "/logout",
		Summary:       "Logout",
		Tags:          []string{"auth"},
		DefaultStatus: http.StatusNoContent,
	}, h.logout)

	if h.Cfg.AppEnv == "development" {
		huma.Register(authAPI, huma.Operation{
			OperationID: "devLogin",
			Method:      http.MethodPost,
			Path:        "/dev/login",
			Summary:     "Dev login",
			Description: "Login into an admin account without credentials in development",
			Tags:        []string{"auth"},
		}, h.devLogin)
	}
}
