package server

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/jackc/pgx/v5"
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

func streamIpInterceptor(queries *database.Queries, config config.GrpcConfig) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		conn, err := pgx.Connect(ctx, database.GetUrl(&config.DB))
		if err != nil {
			logging.Error("failed to get database connection for transaction", err)
			return status.Error(codes.Internal, "")
		}

		tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
		if err != nil {
			logging.Error("failed to start transaction", err)
			return status.Error(codes.Internal, "")
		}
		defer tx.Rollback(ctx)
		queries = queries.WithTx(tx)

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

		err = register(ctx, queries, client)
		if err != nil {
			logging.Error(fmt.Sprintf("failed to register client %s", client.Ip), err)
		}

		err = tx.Commit(ctx)
		if err != nil {
			logging.Error("failed to commit transaction", err)
			return status.Error(codes.Internal, "failed to commit transaction")
		}

		ctxWithClient := withClient(ctx, client)

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
