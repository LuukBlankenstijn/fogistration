package service

import (
	"time"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceRepo struct {
	Auth        *authService
	Client      *clientService
	Team        *teamService
	Wallpaper   *wallpaperService
	OIDCService *oidcService
}

func New(pool *pgxpool.Pool, cfg *config.HttpConfig) *ServiceRepo {
	q := database.New(pool)
	signer := &JWTSigner{
		Secret: []byte(cfg.Secret),
		TTL:    time.Hour,
	}
	authService := newAuthService(q, signer)
	return &ServiceRepo{
		Auth:        authService,
		Client:      newClientService(pool),
		Team:        newTeamService(pool),
		Wallpaper:   newWallpaperService(cfg, pool),
		OIDCService: newOIDCService(cfg.OIDC, authService, pool),
	}
}
