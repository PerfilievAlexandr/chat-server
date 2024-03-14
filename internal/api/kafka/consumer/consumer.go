package kafkaProducer

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/PerfilievAlexandr/chat-server/internal/logger"
	"go.uber.org/zap"
)

type KafkaConsumer struct {
	consumer sarama.Consumer
}

func New(_ context.Context, connectionString string) *KafkaConsumer {
	consumer, err := sarama.NewConsumer([]string{connectionString}, nil)
	if err != nil {
		logger.Fatal("failed create kafka consumer", zap.Any("err", err))
		return nil
	}

	return &KafkaConsumer{consumer}
}
