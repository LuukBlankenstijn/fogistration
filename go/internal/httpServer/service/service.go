package service

import (
	"time"

	"github.com/LuukBlankenstijn/fogistration/internal/httpServer/service/auth"
	"github.com/LuukBlankenstijn/fogistration/internal/httpServer/service/client"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
)

type ServiceRepo struct {
	Auth   auth.Service
	Client client.Service
}

func New(q *database.Queries, secret string) *ServiceRepo {
	signer := &auth.JWTSigner{
		Secret: []byte(secret),
		TTL:    time.Hour,
	}
	return &ServiceRepo{
		Auth:   auth.New(q, signer),
		Client: client.New(q),
	}
}
