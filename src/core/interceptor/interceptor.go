package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
	"vpn/src/core/mapError"
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
				return nil, mapError.MapError(NewMissingMetadata(nil))
			}

			authHeader := md.Get("authorization")
			if len(authHeader) == 0 || !strings.HasPrefix(authHeader[0], "Bearer ") {
				return nil, mapError.MapError(NewInvalidAuthHeader(nil))
			}

			tokenReceived := strings.TrimPrefix(authHeader[0], "Bearer ")
			if err := h.Verify(tokenReceived); err != nil {
				return nil, mapError.MapError(NewInvalidToken(err))
			}
		}
		return handler(ctx, req)
	}
}
