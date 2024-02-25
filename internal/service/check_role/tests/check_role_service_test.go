package tests

import (
	"context"
	checkRoleService "github.com/PerfilievAlexandr/chat-server/internal/service/check_role"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestCheckRole(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		role string
	}

	var (
		ctx        = context.Background()
		serviceErr = status.Errorf(codes.PermissionDenied, "incorrect role")
	)

	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "is admin case",
			args: args{
				ctx:  ctx,
				role: "admin",
			},
			err: nil,
		},
		{
			name: "isn't admin case",
			args: args{
				ctx:  ctx,
				role: "",
			},
			err: serviceErr,
		},
		{
			name: "isn't admin case",
			args: args{
				ctx:  ctx,
				role: "user",
			},
			err: serviceErr,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			checkRoleServiceTest := checkRoleService.NewCheckRoleService(ctx)

			res := checkRoleServiceTest.CheckAdmin(tt.args.ctx, tt.args.role)
			require.Equal(t, tt.err, res)
		})
	}
}
