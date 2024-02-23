package historyDtoDb

import (
	messageStatus "github.com/PerfilievAlexandr/chat-server/internal/domain/enum"
	"github.com/google/uuid"
	"time"
)

type HistoryDb struct {
	Id        uuid.UUID                   `db:"id"`
	Text      string                      `db:"text"`
	Status    messageStatus.MessageStatus `db:"status"`
	MessageId uuid.UUID                   `db:"message_id"`
	CreatedAt time.Time                   `db:"created_at"`
}
