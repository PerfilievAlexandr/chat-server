package domain

import "time"

type Message struct {
	Id        int64
	Text      string
	From      string
	CreatedAt time.Time
}
