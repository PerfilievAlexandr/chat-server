package service

import (
	"context"
	"github.com/PerfilievAlexandr/chat-server/internal/api/grpc/chat/dto"
	"github.com/PerfilievAlexandr/chat-server/internal/domain"
	proto "github.com/PerfilievAlexandr/chat-server/pkg/chat_v1"
	"github.com/google/uuid"
)

type ChatService interface {
	SendMessage(ctx context.Context, req dto.SendMessageRequest) error
	CreateChat(ctx context.Context, req dto.CreateChatRequest) (uuid.UUID, error)
	ConnectChat(req dto.ConnectChatRequest, stream proto.ChatV1_ConnectChatServer) error
}

type CheckRoleService interface {
	CheckAdmin(ctx context.Context, role string) error
}

type HistoryService interface {
	SaveHistory(ctx context.Context, message domain.Message) error
}

type MessageService interface {
	SaveMessage(ctx context.Context, req dto.SendMessageRequest) (domain.Message, error)
	GetMessagesByChatId(ctx context.Context, chatId string) ([]domain.Message, error)
}
