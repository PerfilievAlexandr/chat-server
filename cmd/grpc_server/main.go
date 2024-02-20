package main

import (
	"context"
	"github.com/PerfilievAlexandr/chat-server/internal/app"
	"github.com/PerfilievAlexandr/chat-server/internal/logger"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		logger.Fatal("failed to init app", zap.Any("err", err))
	}

	err = a.Run(ctx)
	if err != nil {
		logger.Fatal("failed to run app", zap.Any("err", err))
	}
}
