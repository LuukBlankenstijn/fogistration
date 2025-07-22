package server

import (
	"context"
	"net"
	"strings"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w wrappedServerStream) Context() context.Context {
	return w.ctx
}

func streamIpInterceptor(queries *database.Queries) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return status.Error(codes.InvalidArgument, "missing metadata")
		}

		ips := md.Get("client-ip")
		if len(ips) == 0 {
			return status.Error(codes.InvalidArgument, "missing client ip")
		}
		ip := ips[0]
		if !isValidIPv4(ip) {
			return status.Error(codes.InvalidArgument, "invalid ipv4")
		}

		client, err := queries.UpsertClient(ss.Context(), ip)
		if err != nil {
			logging.Error("failed to upsert client in database", err)
			return status.Error(codes.Internal, "failed to upsert client")
		}

		ctxWithClient := withClient(ss.Context(), client)

		wrappedServerStream := wrappedServerStream{
			ServerStream: ss,
			ctx:          ctxWithClient,
		}
		return handler(srv, wrappedServerStream)
	}
}

func isValidIPv4(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	return parsedIP.To4() != nil && !strings.Contains(ip, ":")
}
