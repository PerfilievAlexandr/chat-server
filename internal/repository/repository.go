package repository

import (
	"context"
	"github.com/PerfilievAlexandr/chat-server/internal/api/grpc/chat/dto"
	"github.com/PerfilievAlexandr/chat-server/internal/domain"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MessageRepository interface {
	Delete(ctx context.Context, messageId int64) (emptypb.Empty, error)
	SaveMessage(ctx context.Context, req dto.SendMessageRequest) (domain.Message, error)
	GetMessagesByChatId(ctx context.Context, chatId uuid.UUID) ([]domain.Message, error)
}

type ChatRepository interface {
	SaveChat(ctx context.Context, req dto.CreateChatRequest) (uuid.UUID, error)
	IsExists(ctx context.Context, chatId uuid.UUID) (bool, error)
}
