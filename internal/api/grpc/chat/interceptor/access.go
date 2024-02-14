package interceptor

import (
	"context"
	"errors"
	authProto "github.com/PerfilievAlexandr/auth/pkg/access_v1"
	authClient "github.com/PerfilievAlexandr/chat-server/internal/integration/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func AccessInterceptor(authClient authClient.AuthServiceClient) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("metadata is not provided")
		}
		outgoingContext := metadata.NewOutgoingContext(ctx, md)

		_, err := authClient.Check(outgoingContext, &authProto.CheckRequest{EndpointAddress: ""})
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}
