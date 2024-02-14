package dtoDb

import (
	"github.com/google/uuid"
	"time"
)

type MessageDb struct {
	Id        uuid.UUID `db:"id"`
	Text      string    `db:"text"`
	From      string    `db:"owner"`
	ChatId    uuid.UUID `db:"chat_id"`
	CreatedAt time.Time `db:"created_at"`
}
