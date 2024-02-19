package checkRole

import (
	"context"
	"github.com/PerfilievAlexandr/chat-server/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	admin = "admin"
)

type checkRoleService struct{}

func NewCheckRoleService(_ context.Context) service.CheckRoleService {
	return &checkRoleService{}
}

func (c *checkRoleService) CheckAdmin(_ context.Context, role string) error {
	if role != admin {
		return status.Errorf(codes.PermissionDenied, "incorrect role")
	}

	return nil
}
