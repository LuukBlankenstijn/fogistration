package container

import (
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/http/sse"
	"github.com/LuukBlankenstijn/fogistration/internal/http-server/service"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	Cfg *config.HttpConfig
	Q   *database.Queries
	S   *service.ServiceRepo
	SSE *sse.SSEManager
}

func NewContainer(cfg *config.HttpConfig, pool *pgxpool.Pool) *Container {
	q := database.New(pool)
	return &Container{
		Cfg: cfg,
		Q:   q,
		S:   service.New(pool, cfg),
		SSE: sse.New(),
	}
}
