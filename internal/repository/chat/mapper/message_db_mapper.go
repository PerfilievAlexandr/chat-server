package mapper

import (
	"github.com/PerfilievAlexandr/chat-server/internal/domain"
	"github.com/PerfilievAlexandr/chat-server/internal/repository/message/dtoDb"
)

func ToMessageFromDbMessage(dbMessage *dtoDb.MessageDb) *domain.Message {
	return &domain.Message{
		Id:        dbMessage.Id,
		Text:      dbMessage.Text,
		From:      dbMessage.From,
		CreatedAt: dbMessage.CreatedAt,
	}
}
