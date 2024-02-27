package chat

import (
	"context"
	"github.com/PerfilievAlexandr/chat-server/internal/api/grpc/chat/dto"
	"github.com/PerfilievAlexandr/chat-server/internal/repository"
	"github.com/PerfilievAlexandr/chat-server/internal/repository/mocks"
	"github.com/PerfilievAlexandr/chat-server/internal/service"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestCreateChat(t *testing.T) {
	t.Parallel()
	type messageServiceMockFunc func(mc *minimock.Controller) service.MessageService
	type chatRepositoryMockFunc func(mc *minimock.Controller) repository.ChatRepository

	type args struct {
		ctx context.Context
		req dto.CreateChatRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		createChatErr = status.Errorf(codes.Internal, "error create chat")

		req = dto.CreateChatRequest{
			Username: "Bob",
		}

		result = uuid.New()
	)

	t.Cleanup(mc.Finish)

	tests := []struct {
		name                       string
		args                       args
		result                     uuid.UUID
		err                        error
		isChannelWithChatIdCreated bool
		chatRepositoryMockFunc     chatRepositoryMockFunc
		messageServiceMockFunc     messageServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err:                        nil,
			result:                     result,
			isChannelWithChatIdCreated: true,
			chatRepositoryMockFunc: func(mc *minimock.Controller) repository.ChatRepository {
				mock := mocks.NewChatRepositoryMock(mc)
				mock.SaveChatMock.Expect(ctx, req).Return(result, nil)
				return mock
			},
			messageServiceMockFunc: func(mc *minimock.Controller) service.MessageService { return nil },
		},
		{
			name: "error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err:                        createChatErr,
			result:                     uuid.UUID{},
			isChannelWithChatIdCreated: false,
			chatRepositoryMockFunc: func(mc *minimock.Controller) repository.ChatRepository {
				mock := mocks.NewChatRepositoryMock(mc)
				mock.SaveChatMock.Expect(ctx, req).Return(uuid.UUID{}, createChatErr)
				return mock
			},
			messageServiceMockFunc: func(mc *minimock.Controller) service.MessageService { return nil },
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			messageServiceMock := tt.messageServiceMockFunc(mc)
			chatRepositoryMock := tt.chatRepositoryMockFunc(mc)
			serviceChatTest := NewChatService(ctx, messageServiceMock, chatRepositoryMock).(*chatService)

			res, err := serviceChatTest.CreateChat(tt.args.ctx, tt.args.req)
			_, ok := serviceChatTest.channels[res.String()]

			require.Equal(t, tt.err, err)
			require.Equal(t, tt.isChannelWithChatIdCreated, ok)
			require.Equal(t, tt.result, res)
		})
	}
}
