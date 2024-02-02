package authClient

import (
	"context"
	authProto "github.com/PerfilievAlexandr/auth/pkg/access_v1"
	authClient "github.com/PerfilievAlexandr/chat-server/internal/integration/auth"
	"google.golang.org/protobuf/types/known/emptypb"
)

type client struct {
	authClient authProto.AccessV1Client
}

func New(authClient authProto.AccessV1Client) authClient.AuthServiceClient {
	return &client{authClient: authClient}
}

func (c *client) Check(ctx context.Context, req *authProto.CheckRequest) (*emptypb.Empty, error) {
	return c.authClient.Check(ctx, &authProto.CheckRequest{EndpointAddress: req.EndpointAddress})
}
