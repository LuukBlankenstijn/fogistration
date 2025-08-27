package handlers

import (
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/container"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/handlers/auth"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/handlers/clients"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/handlers/contests"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/handlers/teams"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/handlers/users"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/handlers/wallpapers"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/middleware"
	"github.com/danielgtaylor/huma/v2"
)

type Handlers struct {
	auth       *auth.Handlers
	user       *users.Handlers
	client     *clients.Handlers
	team       *teams.Handlers
	contests   *contests.Handlers
	wallpapers *wallpapers.Handlers
}

func NewHandlers(container *container.Container) *Handlers {
	return &Handlers{
		auth:       auth.NewHandlers(container),
		user:       users.NewHandlers(container),
		client:     clients.NewHandlers(container),
		team:       teams.NewHandlers(container),
		contests:   contests.NewHandlers(container),
		wallpapers: wallpapers.NewHandlers(container),
	}
}

func (h *Handlers) Register(
	api huma.API,
	middlewareFactory *middleware.MiddlewareFactory,
	prefixes ...string,
) {
	if len(prefixes) > 0 {
		api = huma.NewGroup(api, prefixes...)
	}

	h.auth.Register(api, middlewareFactory, "/auth")
	h.user.Register(api, middlewareFactory, "/users")
	h.client.Register(api, middlewareFactory, "/clients")
	h.team.Register(api, middlewareFactory, "/teams")
	h.contests.Register(api, middlewareFactory, "/contests")
	h.wallpapers.Register(api, middlewareFactory, "/wallpapers")
}
