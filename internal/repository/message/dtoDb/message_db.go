package dtoDb

import (
	messageStatus "github.com/PerfilievAlexandr/chat-server/internal/domain/enum"
	"github.com/google/uuid"
	"time"
)

type MessageDb struct {
	Id        uuid.UUID                   `db:"id"`
	Text      string                      `db:"text"`
	From      string                      `db:"owner"`
	Status    messageStatus.MessageStatus `db:"status"`
	ChatId    uuid.UUID                   `db:"chat_id"`
	CreatedAt time.Time                   `db:"created_at"`
}
