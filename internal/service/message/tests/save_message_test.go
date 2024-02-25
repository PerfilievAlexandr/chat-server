package tests

import (
	"testing"
)

func TestSaveHistory(t *testing.T) {
	//t.Parallel()
	//type historyServiceMockFunc func(mc *minimock.Controller) service.HistoryService
	//type messageRepositoryMockFunc func(mc *minimock.Controller) repository.MessageRepository
	//type txManagerMockFunc func(mc *minimock.Controller) test_utils.TxManager
	//
	//type args struct {
	//	ctx context.Context
	//	req dto.SendMessageRequest
	//}
	//
	//var (
	//	ctx = context.Background()
	//	mc  = minimock.NewController(t)
	//
	//	saveMessageErr = status.Errorf(codes.Internal, "error save message")
	//	//saveHistoryErr = status.Errorf(codes.Internal, "error save history")
	//
	//	req = dto.SendMessageRequest{
	//		ChatId: uuid.New().String(),
	//	}
	//
	//	result = domain.Message{
	//		Id:        uuid.New(),
	//		Text:      "Test",
	//		From:      "user1",
	//		Status:    messageStatus.NEW,
	//		CreatedAt: time.Now(),
	//	}
	//
	//	txFunc = func(ctx context.Context) error {
	//		return nil
	//	}
	//)
	//
	//t.Cleanup(mc.Finish)

	//tests := []struct {
	//	name                      string
	//	args                      args
	//	result                    domain.Message
	//	err                       error
	//	messageRepositoryMockFunc messageRepositoryMockFunc
	//	historyServiceMockFunc    historyServiceMockFunc
	//	txManagerMockFunc         txManagerMockFunc
	//}{
	//	{
	//		name: "success case",
	//		args: args{
	//			ctx: ctx,
	//			req: req,
	//		},
	//		err:    nil,
	//		result: result,
	//		messageRepositoryMockFunc: func(mc *minimock.Controller) repository.MessageRepository {
	//			mock := mocks.NewMessageRepositoryMock(mc)
	//			mock.SaveMessageMock.Expect(ctx, req).Return(result, nil)
	//			return mock
	//		},
	//		historyServiceMockFunc: func(mc *minimock.Controller) service.HistoryService {
	//			mock := serviceMock.NewHistoryServiceMock(mc)
	//			mock.SaveHistoryMock.Expect(ctx, result).Return(nil)
	//			return mock
	//		},
	//		txManagerMockFunc: func(mc *minimock.Controller) test_utils.TxManager {
	//			mock := dbMocks.NewTxManagerMock(mc)
	//			mock.ReadCommittedMock.Expect(ctx, txFunc).Return(nil)
	//			return mock
	//		},
	//	},
	//	{
	//		name: "error case",
	//		args: args{
	//			ctx: ctx,
	//			req: req,
	//		},
	//		err:    saveMessageErr,
	//		result: result,
	//		messageRepositoryMockFunc: func(mc *minimock.Controller) repository.MessageRepository {
	//			mock := mocks.NewMessageRepositoryMock(mc)
	//			mock.SaveMessageMock.Expect(ctx, req).Return(domain.Message{}, saveMessageErr)
	//			return mock
	//		},
	//		historyServiceMockFunc: func(mc *minimock.Controller) service.HistoryService {
	//			mock := serviceMock.NewHistoryServiceMock(mc)
	//			mock.SaveHistoryMock.Expect(ctx, result).Return(nil)
	//			return mock
	//		},
	//		txManagerMockFunc: func(mc *minimock.Controller) test_utils.TxManager {
	//			mock := dbMocks.NewTxManagerMock(mc)
	//			mock.ReadCommittedMock.Expect(ctx, txFunc).Return(saveMessageErr)
	//			return mock
	//		},
	//	},
	//}
	//
	//for _, tt := range tests {
	//	tt := tt
	//	t.Run(tt.name, func(t *testing.T) {
	//		t.Parallel()
	//
	//		historyServiceMock := tt.historyServiceMockFunc(mc)
	//		messageRepositoryMock := tt.messageRepositoryMockFunc(mc)
	//		txManagerMock := tt.txManagerMockFunc(mc)
	//		serviceHistoryTest := serviceMessage.NewMessageService(ctx, messageRepositoryMock, historyServiceMock, txManagerMock)
	//
	//		res, err := serviceHistoryTest.SaveMessage(tt.args.ctx, tt.args.req)
	//		require.Equal(t, tt.err, err)
	//		require.Equal(t, tt.result, res)
	//	})
	//}
}
