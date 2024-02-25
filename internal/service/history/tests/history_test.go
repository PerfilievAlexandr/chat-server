package tests

import (
	"context"
	"github.com/PerfilievAlexandr/chat-server/internal/domain"
	"github.com/PerfilievAlexandr/chat-server/internal/repository"
	"github.com/PerfilievAlexandr/chat-server/internal/repository/mocks"
	serviceHistory "github.com/PerfilievAlexandr/chat-server/internal/service/history"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestSaveHistory(t *testing.T) {
	t.Parallel()
	type historyRepositoryMockFunc func(mc *minimock.Controller) repository.HistoryRepository

	type args struct {
		ctx     context.Context
		message domain.Message
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		serviceErr = status.Errorf(codes.Internal, "error save history")

		message = domain.Message{
			Id: uuid.New(),
		}
	)

	t.Cleanup(mc.Finish)

	tests := []struct {
		name                  string
		args                  args
		err                   error
		historyRepositoryMock historyRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:     ctx,
				message: message,
			},
			err: nil,
			historyRepositoryMock: func(mc *minimock.Controller) repository.HistoryRepository {
				mock := mocks.NewHistoryRepositoryMock(mc)
				mock.SaveHistoryMock.Expect(ctx, message).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx:     ctx,
				message: message,
			},
			err: serviceErr,
			historyRepositoryMock: func(mc *minimock.Controller) repository.HistoryRepository {
				mock := mocks.NewHistoryRepositoryMock(mc)
				mock.SaveHistoryMock.Expect(ctx, message).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			historyRepositoryMock := tt.historyRepositoryMock(mc)
			serviceHistoryTest := serviceHistory.NewHistoryService(ctx, historyRepositoryMock)

			err := serviceHistoryTest.SaveHistory(tt.args.ctx, tt.args.message)
			require.Equal(t, tt.err, err)
		})
	}
}
