package service

import (
	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceRepo struct {
	Client    *clientService
	Team      *teamService
	Wallpaper *wallpaperService
}

func New(pool *pgxpool.Pool, cfg *config.HttpConfig) *ServiceRepo {
	return &ServiceRepo{
		Client:    newClientService(pool),
		Team:      newTeamService(pool),
		Wallpaper: newWallpaperService(cfg, pool),
	}
}
