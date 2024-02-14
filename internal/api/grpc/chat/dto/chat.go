package dto

import "time"

type CreateChatRequest struct {
	Username string
}

type SendMessageRequest struct {
	ChatId    string
	Text      string
	Owner     string
	CreatedAt time.Time
}

type ConnectChatRequest struct {
	ChatId   string
	Username string
}
