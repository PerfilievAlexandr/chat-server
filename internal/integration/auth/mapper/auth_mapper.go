package mapper

import (
	authProto "github.com/PerfilievAlexandr/auth/pkg/access_v1"
	"github.com/PerfilievAlexandr/chat-server/internal/integration/auth/dto"
)

func MapToClaimsResponse(claims *authProto.ClaimsResponse) dto.ClaimsResponse {
	return dto.ClaimsResponse{
		Username: claims.Username,
		Role:     claims.Role,
	}
}
