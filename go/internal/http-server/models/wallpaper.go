package models

import (
	"time"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database/models"
)

type Wallpaper struct {
	ID        int32                   `json:"id"`
	Layout    *models.WallpaperLayout `json:"layout,omitempty"`
	CreatedAt time.Time               `json:"created_at" readOnly:"true"`
	UpdatedAt time.Time               `json:"updated_at" readOnly:"true"`
}

func MapWallpaper(wallpapers ...database.Wallpaper) []Wallpaper {
	newWallpapers := []Wallpaper{}
	for _, wallpaper := range wallpapers {
		newWallpaper := Wallpaper{
			ID:        wallpaper.ID,
			Layout:    wallpaper.Layout,
			CreatedAt: wallpaper.CreatedAt.Time,
			UpdatedAt: wallpaper.UpdatedAt.Time,
		}
		newWallpapers = append(newWallpapers, newWallpaper)
	}
	return newWallpapers
}

type WallpaperLayout = models.WallpaperLayout
