package authClient

import (
	"context"
	authProto "github.com/PerfilievAlexandr/auth/pkg/access_v1"
	authClient "github.com/PerfilievAlexandr/chat-server/internal/integration/auth"
	"github.com/PerfilievAlexandr/chat-server/internal/integration/auth/dto"
	"github.com/PerfilievAlexandr/chat-server/internal/integration/auth/mapper"
	"google.golang.org/protobuf/types/known/emptypb"
)

type client struct {
	authClient authProto.AccessV1Client
}

func New(authClient authProto.AccessV1Client) authClient.AuthServiceClient {
	return &client{authClient: authClient}
}

func (c *client) Check(ctx context.Context) (dto.ClaimsResponse, error) {
	claims, err := c.authClient.Check(ctx, &emptypb.Empty{})
	if err != nil {
		return dto.ClaimsResponse{}, err
	}

	return mapper.MapToClaimsResponse(claims), nil
}
