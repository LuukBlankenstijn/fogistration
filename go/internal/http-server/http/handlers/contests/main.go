package contests

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
	groupApi := huma.NewGroup(api, prefixes...)
	groupApi.UseMiddleware(middlewareFactory.Auth(groupApi)...)

	huma.Register(groupApi, huma.Operation{
		OperationID: "getActiveContest",
		Method:      http.MethodGet,
		Path:        "/active",
		Tags:        []string{"contests"},
	}, h.getActiveContest)

	huma.Register(groupApi, huma.Operation{
		OperationID: "listContests",
		Method:      http.MethodGet,
		Path:        "/",
		Tags:        []string{"contests"},
	}, h.getAllContests)
}
