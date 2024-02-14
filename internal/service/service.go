package service

import (
	"context"
	"github.com/PerfilievAlexandr/chat-server/internal/api/grpc/chat/dto"
	proto "github.com/PerfilievAlexandr/chat-server/pkg/chat_v1"
	"github.com/google/uuid"
)

type ChatService interface {
	SendMessage(ctx context.Context, req dto.SendMessageRequest) error
	CreateChat(ctx context.Context, req dto.CreateChatRequest) (uuid.UUID, error)
	ConnectChat(req dto.ConnectChatRequest, stream proto.ChatV1_ConnectChatServer) error
}
