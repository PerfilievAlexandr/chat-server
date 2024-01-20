package mapper

import (
	"chat-server/internal/api/grpc/message/dto"
	proto "chat-server/pkg/chat_v1"
)

func MapToCreateRequest(req *proto.CreteRequest) *dto.CreateRequest {
	return &dto.CreateRequest{
		Usernames: req.Usernames,
	}
}

func MapToSendMessageRequest(req *proto.SendMessageRequest) *dto.SendMessageRequest {
	return &dto.SendMessageRequest{
		Text:      req.Text,
		From:      req.From,
		CreatedAt: req.Timestamp.AsTime(),
	}
}
