package mapper

import (
	"github.com/PerfilievAlexandr/chat-server/internal/domain"
	"github.com/PerfilievAlexandr/chat-server/internal/repository/message/dtoDb"
)

func ToMessageFromDbMessage(dbMessage dtoDb.MessageDb) domain.Message {
	return domain.Message{
		Id:        dbMessage.Id,
		Text:      dbMessage.Text,
		From:      dbMessage.From,
		Status:    dbMessage.Status,
		CreatedAt: dbMessage.CreatedAt,
	}
}

func ToMessagesFromDbMessages(dbMessage []dtoDb.MessageDb) []domain.Message {
	var messages []domain.Message
	for _, message := range dbMessage {
		res := ToMessageFromDbMessage(message)
		messages = append(messages, res)
	}

	return messages
}
