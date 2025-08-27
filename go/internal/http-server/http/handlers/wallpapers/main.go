package wallpapers

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
		OperationID: "getWallpaperLayout",
		Method:      http.MethodGet,
		Path:        "/{id}/layout",
		Summary:     "Get wallpaper layout",
		Tags:        []string{"wallpapers"},
	}, h.getWallpaperLayout)

	huma.Register(groupApi, huma.Operation{
		OperationID: "putWallpaperLayout",
		Method:      http.MethodPut,
		Path:        "/{id}/layout",
		Summary:     "Creates or replaces wallpaper with given layout",
		Tags:        []string{"wallpapers"},
	}, h.putWallpaperLayout)

	huma.Register(groupApi, huma.Operation{
		OperationID: "deleteWallpaperLayout",
		Method:      http.MethodDelete,
		Path:        "/{id}/layout",
		Summary:     "Delete wallpaper layout",
		Tags:        []string{"wallpapers"},
	}, h.deleteWallpaperLayout)

	huma.Register(groupApi, huma.Operation{
		OperationID: "getWallpaperFile",
		Method:      http.MethodGet,
		Path:        "/{id}/file",
		Summary:     "Get wallpaper file by id",
		Tags:        []string{"wallpapers"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Image response",
				Content: map[string]*huma.MediaType{
					"image/png": {Schema: &huma.Schema{
						Type:   "string",
						Format: "binary",
					}},
				},
			},
		},
	}, h.getWallpaperFile)

	huma.Register(groupApi, huma.Operation{
		OperationID: "putWallpaperFile",
		Method:      http.MethodPut,
		Path:        "/{id}/file",
		Summary:     "Create or replace wallpaper file",
		Tags:        []string{"wallpapers"},
		RequestBody: &huma.RequestBody{
			Content: map[string]*huma.MediaType{
				"image/png": {},
			},
		},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Image response",
				Content: map[string]*huma.MediaType{
					"image/png": {Schema: &huma.Schema{
						Type:   "string",
						Format: "binary",
					}},
				},
			},
		},
	}, h.putWallpaperFile)

	huma.Register(groupApi, huma.Operation{
		OperationID: "deleteWallpaperFile",
		Method:      http.MethodDelete,
		Path:        "/{id}/file",
		Summary:     "Delete wallpaper file",
		Description: "If the file does not exist, returns succesfully",
		Tags:        []string{"wallpapers"},
	}, h.deleteWallpaperFile)
}
