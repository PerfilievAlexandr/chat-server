package config

import (
	"fmt"
	configInterface "github.com/PerfilievAlexandr/chat-server/internal/config/interface"
	"os"

	"github.com/pkg/errors"
)

var _ configInterface.KafkaConfig = (*kafkaConfig)(nil)

const (
	port = "KAFKA_PORT"
	host = "KAFKA_HOST"
)

type kafkaConfig struct {
	port string
	host string
}

func New() (configInterface.KafkaConfig, error) {
	kafkaPort := os.Getenv(port)
	if len(kafkaPort) == 0 {
		return nil, errors.New("kafka port not found")
	}

	kafkaHost := os.Getenv(host)
	if len(kafkaHost) == 0 {
		return nil, errors.New("kafka host not found")
	}

	return &kafkaConfig{
		host: kafkaHost,
		port: kafkaPort,
	}, nil
}

func (s *kafkaConfig) ConnectString() string {
	return fmt.Sprintf(
		"%s:%s",
		s.host,
		s.port,
	)
}
