package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
	"vpn/src/core/settings"
	"vpn/src/infra/hasher"
)

func AuthUnaryInterceptor(h *hasher.Hasher) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if info.FullMethod == "/config.ConfigService/StartServer" {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
			}

			authHeader := md.Get("authorization")
			if len(authHeader) == 0 || !strings.HasPrefix(authHeader[0], "Bearer ") {
				return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
			}

			tokenReceived := strings.TrimPrefix(authHeader[0], "Bearer ")
			if err := h.Verify(tokenReceived, settings.Config.HashPass); err != nil {
				return nil, status.Errorf(codes.Unauthenticated, "invalid token")
			}
		}
		return handler(ctx, req)
	}
}
