package message

import (
	"context"
	"fmt"
	"github.com/PerfilievAlexandr/chat-server/internal/api/grpc/chat/dto"
	"github.com/PerfilievAlexandr/chat-server/internal/domain"
	"github.com/PerfilievAlexandr/chat-server/internal/repository/message/dtoDb"
	"github.com/PerfilievAlexandr/chat-server/internal/repository/message/mapper"
	"github.com/google/uuid"

	"github.com/PerfilievAlexandr/chat-server/internal/repository"
	"github.com/PerfilievAlexandr/platform_common/pkg/db"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type messageRepository struct {
	db db.Client
}

func NewMessageRepo(_ context.Context, db db.Client) repository.MessageRepository {
	return &messageRepository{db: db}
}

func (m *messageRepository) Delete(ctx context.Context, messageId int64) (emptypb.Empty, error) {
	query := fmt.Sprintf("DELETE FROM messages WHERE id=$1")
	_, err := m.db.ExecContext(ctx, query, messageId)
	if err != nil {
		return emptypb.Empty{}, err
	}

	return emptypb.Empty{}, nil
}

func (m *messageRepository) SaveMessage(ctx context.Context, req dto.SendMessageRequest) (domain.Message, error) {
	id, idErr := uuid.NewUUID()
	if idErr != nil {
		return domain.Message{}, idErr
	}

	query := fmt.Sprintf("INSERT INTO messages (id, text, owner, chat_id, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, text, owner, chat_id, created_at")

	var message = dtoDb.MessageDb{}
	err := m.db.ScanOneContext(ctx, &message, query, id, req.Text, req.Owner, uuid.MustParse(req.ChatId), time.Now())
	if err != nil {
		return domain.Message{}, err
	}

	return mapper.ToMessageFromDbMessage(message), nil
}

func (m *messageRepository) GetMessagesByChatId(ctx context.Context, chatId uuid.UUID) ([]domain.Message, error) {
	query := fmt.Sprintf("SELECT * FROM messages WHERE chat_id = $1")

	var messages []dtoDb.MessageDb
	err := m.db.ScanAllContext(ctx, &messages, query, chatId)
	if err != nil {
		return nil, err
	}

	return mapper.ToMessagesFromDbMessages(messages), err
}
