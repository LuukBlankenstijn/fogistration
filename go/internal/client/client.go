package client

import (
	"context"
	"fmt"
	"reflect"

	"github.com/LuukBlankenstijn/fogistration/internal/client/service"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	client   pb.FogistrationServiceClient
	handlers map[reflect.Type]MessageHandler
	config   config.ClientConfig
}

func NewClient(config config.ClientConfig) (*Client, error) {
	conn, err := grpc.NewClient(config.Server, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewFogistrationServiceClient(conn)

	return &Client{
		client:   client,
		handlers: make(map[reflect.Type]MessageHandler),
		config:   config,
	}, nil
}

func (c *Client) RegisterHandler(handler MessageHandler) {
	handler.SetConfig(c.config)
	c.handlers[handler.MessageType()] = handler
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

		messageType := reflect.TypeOf(msg.GetMessage())
		if handler, ok := c.handlers[messageType]; ok {
			handler.HandleMessage(msg)
		} else {
			logging.Warn("no handler registered for message type %s", messageType)
			logging.Warn("failed to handle message %+v", msg.GetMessage())
		}
	}
}
