package mapper

import (
	"github.com/PerfilievAlexandr/chat-server/internal/domain"
	"github.com/PerfilievAlexandr/chat-server/internal/repository/message/dto"
)

func ToMessageFromDbMessage(dbMessage *dto.MessageDb) *domain.Message {
	return &domain.Message{
		Id:        dbMessage.Id,
		Text:      dbMessage.Text,
		From:      dbMessage.From,
		CreatedAt: dbMessage.CreatedAt,
	}
}
