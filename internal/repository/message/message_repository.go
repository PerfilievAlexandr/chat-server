package message

import (
	"chat-server/internal/api/grpc/message/dto"
	"chat-server/internal/client/db"
	"chat-server/internal/repository"
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type messageRepository struct {
	db db.Client
}

func NewMessageRepo(ctx context.Context, db db.Client) repository.MessageRepository {
	return messageRepository{db: db}
}

func (m messageRepository) Create(ctx context.Context, usernames []string) (int64, error) {
	// TODO
	return 1, nil
}

func (m messageRepository) Delete(ctx context.Context, messageId int64) (*emptypb.Empty, error) {
	query := fmt.Sprintf("DELETE FROM messages WHERE id=$1")
	_, err := m.db.ExecContext(ctx, query, messageId)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (m messageRepository) SendMessage(ctx context.Context, req *dto.SendMessageRequest) (*emptypb.Empty, error) {
	query := fmt.Sprintf("INSERT INTO messages (text, producer, created_at) VALUES ($1, $2, $3)")
	_, err := m.db.ExecContext(ctx, query, req.Text, req.From, time.Now())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
