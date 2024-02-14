package mapper

import (
	"github.com/PerfilievAlexandr/chat-server/internal/api/grpc/chat/dto"
	"github.com/PerfilievAlexandr/chat-server/internal/domain"
	proto "github.com/PerfilievAlexandr/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapToProtoMessage(message *domain.Message) *proto.Message {
	return &proto.Message{
		From:      message.From,
		Text:      message.Text,
		Timestamp: timestamppb.New(message.CreatedAt),
	}
}

func MapToSendMessageRequest(req *proto.SendMessageRequest) dto.SendMessageRequest {
	return dto.SendMessageRequest{
		ChatId:    req.ChatId,
		Text:      req.Message.Text,
		Owner:     req.Message.From,
		CreatedAt: req.Message.Timestamp.AsTime(),
	}
}

func MapToConnectChatRequest(req *proto.ConnectChatRequest) dto.ConnectChatRequest {
	return dto.ConnectChatRequest{
		ChatId:   req.ChatId,
		Username: req.Username,
	}
}

func MapToCreateChatRequest(req *proto.CreateChatRequest) dto.CreateChatRequest {
	return dto.CreateChatRequest{
		Username: req.Username,
	}
}
