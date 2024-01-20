package service

import (
	"chat-server/internal/api/grpc/message/dto"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MessageService interface {
	Create(ctx context.Context, usernames []string) (int64, error)
	Delete(ctx context.Context, messageId int64) (*emptypb.Empty, error)
	SendMessage(ctx context.Context, req *dto.SendMessageRequest) (*emptypb.Empty, error)
}
