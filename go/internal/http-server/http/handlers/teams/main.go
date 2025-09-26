package teams

import (
	"net/http"

	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/container"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/middleware"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/sse"
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
	huma.Register(clientApi, huma.Operation{
		OperationID: "getPrintInfo",
		Method:      http.MethodGet,
		Path:        "/print",
		Summary:     "Get print info for a team",
		Tags:        []string{"teams"},
	}, h.getPrintInfo)

	clientApi.UseMiddleware(middlewareFactory.Auth(clientApi)...)

	sse.Register(h.SSE, clientApi, huma.Operation{
		OperationID: "getTeam",
		Method:      http.MethodGet,
		Path:        "/{id}",
		Tags:        []string{"teams"},
	}, h.getSingleTeam)

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
