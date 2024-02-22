package domain

import (
	"github.com/PerfilievAlexandr/chat-server/internal/domain/enum"
	"github.com/google/uuid"
	"time"
)

type Message struct {
	Id        uuid.UUID
	Text      string
	From      string
	Status    messageStatus.MessageStatus
	CreatedAt time.Time
}
