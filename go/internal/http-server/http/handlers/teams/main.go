package teams

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
		OperationID: "listTeams",
		Method:      http.MethodGet,
		Path:        "/",
		Tags:        []string{"teams"},
	}, h.listTeams)

	huma.Register(clientApi, huma.Operation{
		OperationID: "setTeamClient",
		Method:      http.MethodPost,
		Path:        "/{id}/client",
		Summary:     "Set team client",
		Tags:        []string{"teams"},
	}, h.setClient)
}
