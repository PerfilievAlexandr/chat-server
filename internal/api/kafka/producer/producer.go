package kafkaProducer

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/PerfilievAlexandr/chat-server/internal/logger"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func New(_ context.Context, connectionString string) *KafkaProducer {
	producer, err := sarama.NewSyncProducer([]string{connectionString}, nil)
	if err != nil {
		logger.Fatal("failed create kafka producer", zap.Any("err", err))
		return nil
	}

	return &KafkaProducer{producer}
}

func (k *KafkaProducer) SendCreateChatMessage(ctx context.Context, chatId uuid.UUID) error {
	topic := "createChat"
	messageBytes, err := json.Marshal(chatId)
	if err != nil {
		return status.Errorf(codes.Internal, "kafka producer: error to marshal message. %s", err)
	}

	return k.sendMessage(ctx, messageBytes, topic)
}

func (k *KafkaProducer) sendMessage(_ context.Context, message interface{}, topic string) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return status.Errorf(codes.Internal, "kafka producer: error to marshal message. %s", err)
	}

	msgId := sarama.RecordHeader{
		Key:   []byte("messageId"),
		Value: []byte(uuid.New().String()),
	}
	kafkaMessage := &sarama.ProducerMessage{
		Topic:   topic,
		Value:   sarama.ByteEncoder(messageBytes),
		Headers: []sarama.RecordHeader{msgId},
	}

	_, _, err = k.producer.SendMessage(kafkaMessage)
	if err != nil {
		return status.Errorf(codes.Internal, "kafka producer: error send message. %s", err)
	}

	return nil
}
