package authClient

import (
	"context"
	authProto "github.com/PerfilievAlexandr/auth/pkg/access_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthServiceClient interface {
	Check(ctx context.Context, req *authProto.CheckRequest) (*emptypb.Empty, error)
}
