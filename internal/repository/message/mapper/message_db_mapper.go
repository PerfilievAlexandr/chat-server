package mapper

import (
	"chat-server/internal/domain"
	"chat-server/internal/repository/message/dto"
)

func ToMessageFromDbMessage(dbMessage *dto.MessageDb) *domain.Message {
	return &domain.Message{
		Id:        dbMessage.Id,
		Text:      dbMessage.Text,
		From:      dbMessage.From,
		CreatedAt: dbMessage.CreatedAt,
	}
}
