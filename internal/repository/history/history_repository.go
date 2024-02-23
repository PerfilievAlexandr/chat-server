package historyRepository

import (
	"context"
	"fmt"
	"github.com/PerfilievAlexandr/chat-server/internal/domain"
	"github.com/google/uuid"

	"github.com/PerfilievAlexandr/chat-server/internal/repository"
	"github.com/PerfilievAlexandr/platform_common/pkg/db"
	"time"
)

type historyRepository struct {
	db db.Client
}

func NewHistoryRepo(_ context.Context, db db.Client) repository.HistoryRepository {
	return &historyRepository{db: db}
}

func (m *historyRepository) SaveHistory(ctx context.Context, message domain.Message) error {
	query := fmt.Sprintf("INSERT INTO history (id, text, status, message_id, created_at) VALUES ($1, $2, $3, $4, $5)")

	_, err := m.db.ExecContext(ctx, query, uuid.New(), message.Text, message.Status, message.Id, time.Now())
	if err != nil {
		return err
	}

	return nil
}
