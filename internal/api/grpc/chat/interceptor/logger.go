package interceptor

import (
	"context"
	"github.com/PerfilievAlexandr/chat-server/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	res, err := handler(ctx, req)
	if err != nil {
		logger.Error(err.Error(), zap.String("method", info.FullMethod), zap.Any("req", req))
	}

	return res, err
}
