package message

import (
	"context"
	"github.com/PerfilievAlexandr/chat-server/internal/api/grpc/message/dto"
	"github.com/PerfilievAlexandr/chat-server/internal/repository"
	"github.com/PerfilievAlexandr/chat-server/internal/service"
	"github.com/PerfilievAlexandr/platform_common/pkg/db"
	"google.golang.org/protobuf/types/known/emptypb"
)

type messageService struct {
	repo repository.MessageRepository
	tx   db.TxManager
}

func NewMessageService(ctx context.Context, repo repository.MessageRepository, tx db.TxManager) service.MessageService {
	return &messageService{repo, tx}
}

func (m messageService) Create(ctx context.Context, usernames []string) (int64, error) {
	return m.repo.Create(ctx, usernames)
}

func (m messageService) Delete(ctx context.Context, messageId int64) (*emptypb.Empty, error) {
	return m.repo.Delete(ctx, messageId)
}

func (m messageService) SendMessage(ctx context.Context, req *dto.SendMessageRequest) (*emptypb.Empty, error) {
	return m.repo.SendMessage(ctx, req)
}
