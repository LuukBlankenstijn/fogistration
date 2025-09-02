package server

import (
	"fmt"
	"net"

	"github.com/LuukBlankenstijn/fogistration/internal/grpc/pubsub"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	pb "github.com/LuukBlankenstijn/fogistration/internal/shared/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedFogistrationServiceServer
	pubsub  *pubsub.Manager
	queries *database.Queries
	config  config.GrpcConfig
}

func NewServer(queries *database.Queries, pubsub *pubsub.Manager, config config.GrpcConfig) *Server {
	return &Server{
		pubsub:  pubsub,
		queries: queries,
		config:  config,
	}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", ":"+s.config.Port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", s.config.Port, err)
	}

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(streamIpInterceptor(s.queries, s.config)),
	)
	pb.RegisterFogistrationServiceServer(grpcServer, s)

	logging.Info("Starting gRPC server on port %s", s.config.Port)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve gRPC: %w", err)
	}

	return nil
}

func (s *Server) Stream(clientMessage *pb.ClientMessage, stream grpc.ServerStreamingServer[pb.ServerMessage]) error {
	ctx := stream.Context()
	client, ok := getClient(ctx)
	if !ok {
		logging.Error("failed to get client from stream context", nil)
		return status.Error(codes.NotFound, "client not found on context")
	}
	logging.Info("new client connnection: %s", client.Ip)

	ch := s.pubsub.Subscribe(client.Ip)
	defer s.pubsub.Unsubscribe(client.Ip)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-ch:
			if err := stream.Send(msg); err != nil {
				logging.Error("failed to send: %v", err)
			}
		}
	}
}
