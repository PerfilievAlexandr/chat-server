package domain

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	Id        uuid.UUID
	Text      string
	From      string
	CreatedAt time.Time
}
