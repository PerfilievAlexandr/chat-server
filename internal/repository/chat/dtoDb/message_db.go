package dtoDb

import "time"

type MessageDb struct {
	Id        int64     `db:"id"`
	Text      string    `db:"text"`
	From      string    `db:"producer"`
	CreatedAt time.Time `db:"created_at"`
}
