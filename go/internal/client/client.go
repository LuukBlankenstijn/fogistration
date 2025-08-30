package client

import (
	"context"
	"fmt"

	"github.com/LuukBlankenstijn/fogistration/internal/client/service"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	client   pb.FogistrationServiceClient
	handlers map[string]MessageHandler
}

func NewClient(server string) (*Client, error) {
	conn, err := grpc.NewClient(server, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewFogistrationServiceClient(conn)

	return &Client{
		client:   client,
		handlers: make(map[string]MessageHandler),
	}, nil
}

func (c *Client) RegisterHandler(handler MessageHandler, messageType string) {
	c.handlers[messageType] = handler
}

func (c *Client) StartReceiving(ctx context.Context) error {
	ip, err := service.GetIp()
	if err != nil {
		return fmt.Errorf("failed to get client ip: %w", err)
	}

	md := metadata.New(map[string]string{
		"client-ip": ip.String(),
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	stream, err := c.client.Stream(ctx, &pb.ClientMessage{})
	if err != nil {
		return fmt.Errorf("failed to create stream: %w", err)
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			return fmt.Errorf("failed to receive message: %w", err)
		}

		messageType := getMessageType(msg)
		if handler, ok := c.handlers[messageType]; ok {
			handler.HandleMessage(msg)
		} else {
			logging.Warn("no handler registered for message type %s", messageType)
		}
	}
}

func getMessageType(msg *pb.ServerMessage) string {
	var messageType string
	switch msg.GetMessage().(type) {
	case *pb.ServerMessage_Reload:
		messageType = "reload"
	default:
		messageType = "unknown"
	}

	return messageType
}
