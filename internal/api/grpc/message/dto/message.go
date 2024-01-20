package dto

import "time"

type CreateRequest struct {
	Usernames []string
}

type SendMessageRequest struct {
	Text      string
	From      string
	CreatedAt time.Time
}
