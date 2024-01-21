package config

import (
	"context"
	dbConfig "github.com/PerfilievAlexandr/chat-server/internal/config/db"
	grpcConfig "github.com/PerfilievAlexandr/chat-server/internal/config/grpc"
	configInterface "github.com/PerfilievAlexandr/chat-server/internal/config/interface"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	GRPCConfig configInterface.GrpcServerConfig
	DbConfig   configInterface.DatabaseConfig
}

func NewConfig(ctx context.Context) (*Config, error) {
	dbCfg, err := dbConfig.NewDbConfig()
	if err != nil {
		log.Fatalf("failed to config: %s", err.Error())
	}
	grpcCfg, err := grpcConfig.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to config: %s", err.Error())
	}

	return &Config{
		DbConfig:   dbCfg,
		GRPCConfig: grpcCfg,
	}, nil
}

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}
