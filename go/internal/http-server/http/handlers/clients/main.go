package clients

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
	clientApi := huma.NewGroup(api, prefixes...)
	clientApi.UseMiddleware(middlewareFactory.Auth(clientApi)...)

	huma.Register(clientApi, huma.Operation{
		OperationID: "listClients",
		Method:      http.MethodGet,
		Path:        "/",
		Tags:        []string{"clients"},
	}, h.listClients)

	huma.Register(clientApi, huma.Operation{
		OperationID: "setClientTeam",
		Method:      http.MethodPost,
		Path:        "/{id}/team",
		Summary:     "Set client team",
		Tags:        []string{"clients"},
	}, h.setTeam)
}
