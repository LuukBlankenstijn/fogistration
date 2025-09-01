package client

import (
	"reflect"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/pb"
)

type MessageHandler interface {
	HandleMessage(m *pb.ServerMessage)
	SetConfig(c config.ClientConfig)
	MessageType() reflect.Type
}

type ReloadHandler struct {
	config config.ClientConfig
}

func (s *ReloadHandler) MessageType() reflect.Type {
	return reflect.TypeOf(pb.ServerMessage_Reload{})
}

func (s *ReloadHandler) HandleMessage(m *pb.ServerMessage) {
	msg := m.GetReload()
	if msg == nil {
		return
	}

	logging.Info("reload message")
}

func (s *ReloadHandler) SetConfig(config config.ClientConfig) {
	s.config = config
}
