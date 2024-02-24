package interceptor

import (
	"context"
	prometheusMetrics "github.com/PerfilievAlexandr/chat-server/internal/metrics"
	"google.golang.org/grpc"
	"time"
)

func MetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	prometheusMetrics.IncRequestCounter()

	timeStart := time.Now()
	res, err := handler(ctx, req)
	diffTime := time.Since(timeStart)

	if err != nil {
		prometheusMetrics.IncResponseCounter("error", info.FullMethod)
		prometheusMetrics.HistogramResponseTimeObserve("error", diffTime.Seconds())
	} else {
		prometheusMetrics.IncResponseCounter("success", info.FullMethod)
		prometheusMetrics.HistogramResponseTimeObserve("success", diffTime.Seconds())
	}

	return res, err
}
