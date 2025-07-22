package server

import (
	"context"
	"fmt"

	"github.com/LuukBlankenstijn/fogistration/internal/grpc/pubsub"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/pb"
	"google.golang.org/grpc"
)

type streamHandler struct {
	stream grpc.BidiStreamingServer[pb.ClientMessage, pb.ServerMessage]
	client database.Client
	ctx    context.Context

	queries    *database.Queries
	pubsub     *pubsub.Manager
	regService *registrationService
}

func (h *streamHandler) run() error {
	logging.Info("new client connection: %s", h.client.Ip)

	ready := make(chan struct{})
	go h.handleOutgoing(ready)
	<-ready

	if err := h.regService.register(h.ctx, h.client); err != nil {
		logging.Error("failed registering client: %v", err)
	}

	return h.handleIncoming()
}

func (h *streamHandler) handleOutgoing(ready chan struct{}) {
	ch := h.pubsub.Subscribe(h.client.Ip)
	defer h.pubsub.Unsubscribe(h.client.Ip)
	close(ready)

	for {
		select {
		case <-h.ctx.Done():
			return
		case msg := <-ch:
			if err := h.stream.Send(msg); err != nil {
				logging.Error("failed to send: %v", err)
			}
		}
	}
}

func (h *streamHandler) handleIncoming() error {

	for {
		select {
		case <-h.ctx.Done():
			logging.Info("Context timeout/cancelled")
			return h.ctx.Err()
		default:
			msg, err := h.stream.Recv()
			if err != nil {
				return fmt.Errorf("stream closed: %w", err)
			}
			err = h.queries.UpdateClientLastSeen(h.ctx, h.client.Ip)
			if err != nil {
				logging.Error("failed to update client last seen", err)
			}

			switch message := msg.Message.(type) {
			case *pb.ClientMessage_Heartbeat:
				// nothing to do, already updated Last seen
			case nil:
				logging.Warn("Empty message from client %s", h.client.Ip)
			default:
				logging.Warn("Unknown message type from client %s: %T", h.client.Ip, message)
			}
		}
	}
}
