package config

import (
	"errors"
	configInterface "github.com/PerfilievAlexandr/chat-server/internal/config/interface"
	"net"
	"os"
)

var _ configInterface.GrpcAuthClientConfig = (*authClientConfig)(nil)

const (
	grpcHostEnvName = "GRPC_AUTH_HOST"
	grpcPortEnvName = "GRPC_AUTH_PORT"
)

type authClientConfig struct {
	host string
	port string
}

func NewAuthClientConfig() (configInterface.GrpcAuthClientConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("auth client host not found")
	}

	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("auth client port not found")
	}

	return &authClientConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *authClientConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
