package tracing

import (
	"github.com/uber/jaeger-client-go/config"
	"log"
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
		log.Fatal("failed to init tracing")
	}
}
