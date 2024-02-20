package tracing

import (
	"github.com/PerfilievAlexandr/chat-server/internal/logger"
	"github.com/uber/jaeger-client-go/config"
)

func Init(serviceName string) {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}

	_, err := cfg.InitGlobalTracer(serviceName)
	if err != nil {
		logger.Fatal("failed to init tracing")
	}
}
