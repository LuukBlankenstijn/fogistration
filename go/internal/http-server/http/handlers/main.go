package handlers

import (
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/container"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/handlers/clients"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/handlers/contests"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/handlers/teams"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/handlers/wallpapers"
	"github.com/danielgtaylor/huma/v2"
)

type Handlers struct {
	client     *clients.Handlers
	team       *teams.Handlers
	contests   *contests.Handlers
	wallpapers *wallpapers.Handlers
}

func NewHandlers(container *container.Container) *Handlers {
	return &Handlers{
		client:     clients.NewHandlers(container),
		team:       teams.NewHandlers(container),
		contests:   contests.NewHandlers(container),
		wallpapers: wallpapers.NewHandlers(container),
	}
}

func (h *Handlers) Register(
	api huma.API,
	prefixes ...string,
) {
	if len(prefixes) > 0 {
		api = huma.NewGroup(api, prefixes...)
	}

	h.client.Register(api, "/clients")
	h.team.Register(api, "/teams")
	h.contests.Register(api, "/contests")
	h.wallpapers.Register(api, "/wallpapers")
}
