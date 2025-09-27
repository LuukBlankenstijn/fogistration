package client

import (
	"reflect"

	"github.com/LuukBlankenstijn/fogistration/internal/client/handlers"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/pb"
)

type MessageHandler interface {
	HandleMessage(m *pb.ServerMessage)
	SetConfig(c config.ClientConfig)
	MessageType() reflect.Type
}

type UpdateHandler = handlers.UpdateHandler
