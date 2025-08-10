package client

import (
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/pb"
)

type MessageHandler interface {
	HandleMessage(m *pb.ServerMessage)
}

type ReloadHandler struct{}

func (s ReloadHandler) HandleMessage(m *pb.ServerMessage) {
	msg := m.GetReload()
	if msg == nil {
		return
	}

	logging.Info("reload message")
}
