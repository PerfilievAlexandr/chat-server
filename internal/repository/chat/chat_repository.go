package message

import (
	"context"
	"fmt"
	"github.com/PerfilievAlexandr/chat-server/internal/api/grpc/chat/dto"
	"github.com/google/uuid"

	"github.com/PerfilievAlexandr/chat-server/internal/repository"
	"github.com/PerfilievAlexandr/platform_common/pkg/db"
	"time"
)

type chatRepository struct {
	db db.Client
}

func NewChatRepo(_ context.Context, db db.Client) repository.ChatRepository {
	return &chatRepository{db}
}

func (c *chatRepository) SaveChat(ctx context.Context, req dto.CreateChatRequest) (uuid.UUID, error) {
	id, idErr := uuid.NewUUID()
	if idErr != nil {
		return uuid.UUID{}, idErr
	}
	query := fmt.Sprintf("INSERT INTO chat (id, owner, created_at) VALUES ($1, $2, $3)")

	_, err := c.db.ExecContext(ctx, query, id, req.Username, time.Now())
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func (c *chatRepository) IsExists(ctx context.Context, chatId uuid.UUID) (bool, error) {
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM chat WHERE id=$1)")
	var contains bool
	err := c.db.ScanOneContext(ctx, &contains, query, chatId)
	if err != nil {
		return false, err
	}

	return contains, nil
}
