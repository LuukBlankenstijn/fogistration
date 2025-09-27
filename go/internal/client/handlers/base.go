package handlers

import "github.com/LuukBlankenstijn/fogistration/internal/shared/config"

type baseHandler struct {
	config config.ClientConfig
}

func (b *baseHandler) SetConfig(config config.ClientConfig) {
	b.config = config
}
