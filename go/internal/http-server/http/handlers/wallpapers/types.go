package wallpapers

import (
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/models"
)

type wallpaperLayoutBody struct {
	Body models.WallpaperLayout
}

type wallpaperIdPath struct {
	ID int32 `path:"id" doc:"Wallpaper id"`
}

type wallpaperFileBody struct {
	ContentType string `header:"Content-Type"`
	Body        []byte
}

type getWallpaperLayoutRequest = wallpaperIdPath

type wallpaperLayoutResponse = wallpaperLayoutBody

type putWallpaperLayoutRequest struct {
	wallpaperLayoutBody
	ID int32 `path:"id" doc:"Wallpaper id"`
}

type deleteWallpaperLayoutRequest = wallpaperIdPath

type getWallpaperFileRequest = wallpaperIdPath

type wallpaperFileResponse = wallpaperFileBody

type putWallpaperFileRequest struct {
	RawBody     []byte `contentType:"image/png" json:"omitempty"`
	ID          int32  `path:"id" doc:"Wallpaper id"`
	ContentType string `header:"Content-Type"`
}

type deleteWallpaperFileRequest = wallpaperIdPath
