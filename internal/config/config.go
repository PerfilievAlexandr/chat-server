package config

import (
	"context"
	authClient "github.com/PerfilievAlexandr/chat-server/internal/config/auth_client"
	dbConfig "github.com/PerfilievAlexandr/chat-server/internal/config/db"
	grpcConfig "github.com/PerfilievAlexandr/chat-server/internal/config/grpc"
	configInterface "github.com/PerfilievAlexandr/chat-server/internal/config/interface"
	"github.com/PerfilievAlexandr/chat-server/internal/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	GRPCConfig       configInterface.GrpcServerConfig
	DbConfig         configInterface.DatabaseConfig
	AuthClientConfig configInterface.GrpcAuthClientConfig
}

func NewConfig(_ context.Context) (*Config, error) {
	dbCfg, err := dbConfig.NewDbConfig()
	if err != nil {
		logger.Fatal("failed to config", zap.Any("err", err))
	}
	grpcCfg, err := grpcConfig.NewGRPCConfig()
	if err != nil {
		logger.Fatal("failed to config", zap.Any("err", err))
	}
	authClientCfg, err := authClient.NewAuthClientConfig()
	if err != nil {
		logger.Fatal("failed to config", zap.Any("err", err))
	}

	return &Config{
		DbConfig:         dbCfg,
		GRPCConfig:       grpcCfg,
		AuthClientConfig: authClientCfg,
	}, nil
}

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}
