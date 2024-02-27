package serviceMessage

import (
	"context"
	"github.com/PerfilievAlexandr/chat-server/internal/api/grpc/chat/dto"
	"github.com/PerfilievAlexandr/chat-server/internal/domain"
	"github.com/PerfilievAlexandr/chat-server/internal/repository"
	"github.com/PerfilievAlexandr/chat-server/internal/service"
	"github.com/PerfilievAlexandr/platform_common/pkg/db"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type messageService struct {
	messageRepository repository.MessageRepository
	historyService    service.HistoryService
	txManager         db.TxManager
}

func NewMessageService(
	_ context.Context,
	messageRepository repository.MessageRepository,
	historyService service.HistoryService,
	txManager db.TxManager,
) service.MessageService {
	return &messageService{
		messageRepository: messageRepository,
		historyService:    historyService,
		txManager:         txManager,
	}
}

// TODO как протестировать метод? Как записать message?
func (m *messageService) SaveMessage(ctx context.Context, req dto.SendMessageRequest) (domain.Message, error) {
	var message domain.Message

	err := m.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		message, errTx = m.messageRepository.SaveMessage(ctx, req)
		if errTx != nil {
			return status.Errorf(codes.Internal, "error save message")
		}

		errTx = m.historyService.SaveHistory(ctx, message)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return domain.Message{}, err
	}

	return message, nil
}

func (m *messageService) GetMessagesByChatId(ctx context.Context, chatId string) ([]domain.Message, error) {
	messages, err := m.messageRepository.GetMessagesByChatId(ctx, uuid.MustParse(chatId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "load messages error")
	}

	return messages, nil
}
