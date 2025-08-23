package eventhandler

import "github.com/LuukBlankenstijn/fogistration/internal/grpc/pubsub"

type EventHandler struct {
	pubsub *pubsub.Manager
}

func New(pubsub *pubsub.Manager) *EventHandler {
	return &EventHandler{
		pubsub: pubsub,
	}
}
