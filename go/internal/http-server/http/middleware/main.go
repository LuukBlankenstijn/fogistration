package middleware

import "github.com/LuukBlankenstijn/fogistration/internal/http-server/http/container"

type MiddlewareFactory struct {
	*container.Container
}

func NewMiddlewareFactory(container *container.Container) *MiddlewareFactory {
	return &MiddlewareFactory{container}
}
