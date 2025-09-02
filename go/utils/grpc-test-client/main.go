package main

import (
	"context"
	"log"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	pb "github.com/LuukBlankenstijn/fogistration/internal/shared/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	// Connect to server
	conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			logging.Error("failed to close database connection", err)
		}
	}()

	client := pb.NewFogistrationServiceClient(conn)

	// Add IP to metadata
	md := metadata.New(map[string]string{
		"client-ip": "192.168.1.101",
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Start streaming with metadata context
	stream, err := client.Stream(ctx, &pb.ClientMessage{})
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}

	log.Println("Connected with IP in metadata, listening for responses...")

	// Listen for server messages
	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Fatalf("Failed to receive: %v", err)
		}
		logging.Info("received message %+v", msg.Message)
	}
}
