package tests

import (
	"context"
	"github.com/PerfilievAlexandr/chat-server/internal/domain"
	messageStatus "github.com/PerfilievAlexandr/chat-server/internal/domain/enum"
	"github.com/PerfilievAlexandr/chat-server/internal/repository"
	"github.com/PerfilievAlexandr/chat-server/internal/repository/mocks"
	"github.com/PerfilievAlexandr/chat-server/internal/service"
	serviceMessage "github.com/PerfilievAlexandr/chat-server/internal/service/message"
	"github.com/PerfilievAlexandr/chat-server/internal/test_utils"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
	"time"
)

func TestGetMessagesByChatId(t *testing.T) {
	t.Parallel()
	type historyServiceMockFunc func(mc *minimock.Controller) service.HistoryService
	type messageRepositoryMockFunc func(mc *minimock.Controller) repository.MessageRepository
	type txManagerMockFunc func(mc *minimock.Controller) test_utils.TxManager

	type args struct {
		ctx context.Context
		req string
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		loadMessageErr = status.Errorf(codes.Internal, "load messages error")

		chatId     = uuid.New().String()
		chatIdUuid = uuid.MustParse(chatId)

		result = []domain.Message{
			{
				Id:        uuid.New(),
				Text:      "Test",
				From:      "user1",
				Status:    messageStatus.NEW,
				CreatedAt: time.Now(),
			},
		}
	)

	t.Cleanup(mc.Finish)

	tests := []struct {
		name                      string
		args                      args
		result                    []domain.Message
		err                       error
		messageRepositoryMockFunc messageRepositoryMockFunc
		historyServiceMockFunc    historyServiceMockFunc
		txManagerMockFunc         txManagerMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: chatId,
			},
			err:    nil,
			result: result,
			messageRepositoryMockFunc: func(mc *minimock.Controller) repository.MessageRepository {
				mock := mocks.NewMessageRepositoryMock(mc)
				mock.GetMessagesByChatIdMock.Expect(ctx, chatIdUuid).Return(result, nil)
				return mock
			},
			historyServiceMockFunc: func(mc *minimock.Controller) service.HistoryService { return nil },
			txManagerMockFunc:      func(mc *minimock.Controller) test_utils.TxManager { return nil },
		},
		{
			name: "error case",
			args: args{
				ctx: ctx,
				req: chatId,
			},
			err:    loadMessageErr,
			result: nil,
			messageRepositoryMockFunc: func(mc *minimock.Controller) repository.MessageRepository {
				mock := mocks.NewMessageRepositoryMock(mc)
				mock.GetMessagesByChatIdMock.Expect(ctx, chatIdUuid).Return(nil, loadMessageErr)
				return mock
			},
			historyServiceMockFunc: func(mc *minimock.Controller) service.HistoryService { return nil },
			txManagerMockFunc:      func(mc *minimock.Controller) test_utils.TxManager { return nil },
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			historyServiceMock := tt.historyServiceMockFunc(mc)
			messageRepositoryMock := tt.messageRepositoryMockFunc(mc)
			txManagerMock := tt.txManagerMockFunc(mc)
			serviceHistoryTest := serviceMessage.NewMessageService(ctx, messageRepositoryMock, historyServiceMock, txManagerMock)

			res, err := serviceHistoryTest.GetMessagesByChatId(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.result, res)
		})
	}
}
