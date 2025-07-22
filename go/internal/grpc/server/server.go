package server

import (
	"fmt"
	"net"

	"github.com/LuukBlankenstijn/fogistration/internal/grpc/pubsub"
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
}

func NewServer(queries *database.Queries, pubsub *pubsub.Manager) *Server {
	return &Server{
		pubsub:  pubsub,
		queries: queries,
	}
}

func (s *Server) Start(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", port, err)
	}

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(streamIpInterceptor(s.queries)),
	)
	pb.RegisterFogistrationServiceServer(grpcServer, s)

	logging.Info("Starting gRPC server on port %s", port)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve gRPC: %w", err)
	}

	return nil
}

func (s *Server) Stream(stream grpc.BidiStreamingServer[pb.ClientMessage, pb.ServerMessage]) error {
	ctx := stream.Context()
	client, ok := getClient(ctx)
	if !ok {
		logging.Error("failed to get client from stream context", nil)
		return status.Error(codes.NotFound, "client not found on context")
	}
	logging.Info("new client connnection: %s", client.Ip)

	handler := &streamHandler{
		stream:     stream,
		client:     client,
		ctx:        ctx,
		pubsub:     s.pubsub,
		regService: &registrationService{s.queries, s.pubsub},
	}

	return handler.run()
}
