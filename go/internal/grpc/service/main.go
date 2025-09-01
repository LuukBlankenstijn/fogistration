package service

import (
	"github.com/LuukBlankenstijn/fogistration/internal/grpc/pubsub"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
)

type ServiceContainer struct {
	Reload *reloadService
}

func New(queries *database.Queries, pubsub *pubsub.Manager, config config.GrpcConfig) *ServiceContainer {
	return &ServiceContainer{
		Reload: &reloadService{queries, pubsub, config},
	}
}
