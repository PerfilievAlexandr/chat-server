package authClient

import (
	"context"
	"github.com/PerfilievAlexandr/chat-server/internal/integration/auth/dto"
)

type AuthServiceClient interface {
	Check(ctx context.Context) (dto.ClaimsResponse, error)
}
