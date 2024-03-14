package app

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/PerfilievAlexandr/chat-server/internal/api/grpc/chat/interceptor"
	"github.com/PerfilievAlexandr/chat-server/internal/config"
	"github.com/PerfilievAlexandr/chat-server/internal/domain"
	"github.com/PerfilievAlexandr/chat-server/internal/logger"
	prometheusMetrics "github.com/PerfilievAlexandr/chat-server/internal/metrics"
	"github.com/PerfilievAlexandr/chat-server/internal/tracing"
	proto "github.com/PerfilievAlexandr/chat-server/pkg/chat_v1"
	"github.com/PerfilievAlexandr/platform_common/pkg/closer"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
)

type App struct {
	diProvider    *diProvider
	grpcServer    *grpc.Server
	prometheus    *http.Server
	kafkaConsumer sarama.Consumer
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()

		err := a.runGrpcServer(ctx)
		if err != nil {
			logger.Fatal("failed to run GRPC grpcAuthServer", zap.Any("err", err))
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runPrometheus(ctx)
		if err != nil {
			logger.Fatal("failed to run prometheusServer", zap.Any("err", err))
		}
	}()

	go func() {
		defer wg.Done()

		partConsumer, err := a.runKafkaConsumer(ctx)
		if err != nil {
			logger.Fatal("failed to run kafka consumer", zap.Any("err", err))
		}

		for {
			select {
			// (обработка входящего сообщения и отправка ответа в Kafka)
			case msg, ok := <-partConsumer.Messages():
				if !ok {
					log.Println("Channel closed, exiting")
					return
				}

				// Десериализация входящего сообщения из JSON
				var receivedMessage domain.Message
				err := json.Unmarshal(msg.Value, &receivedMessage)

				if err != nil {
					log.Printf("Error unmarshaling JSON: %v\n", err)
					continue
				}
			}
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initProvider,
		a.initLogger,
		a.initPrometheus,
		a.initTrace,
		a.initGrpcServer,
		a.initKafkaConsumer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initProvider(ctx context.Context) error {
	a.diProvider = newDiProvider()
	a.diProvider.AuthClient(ctx)
	return nil
}

func (a *App) initGrpcServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				interceptor.ServerTracingInterceptor,
				//interceptor.AccessInterceptor(a.diProvider.authClient),
				interceptor.LogInterceptor,
				interceptor.MetricsInterceptor,
			),
		),
	)
	reflection.Register(a.grpcServer)
	proto.RegisterChatV1Server(a.grpcServer, a.diProvider.ChatServer(ctx))

	return nil
}

func (a *App) initLogger(_ context.Context) error {
	stdout := zapcore.AddSync(os.Stdout)
	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	core := zapcore.NewCore(consoleEncoder, stdout, zap.InfoLevel)
	logger.Init(core)

	return nil
}

func (a *App) initKafkaConsumer(ctx context.Context) error {
	consumer, err := sarama.NewConsumer([]string{a.diProvider.Config(ctx).KafkaConfig.ConnectString()}, nil)
	if err != nil {
		logger.Fatal("failed create kafka consumer", zap.Any("err", err))
		return nil
	}
	a.kafkaConsumer = consumer

	return nil
}

func (a *App) runKafkaConsumer(_ context.Context) (sarama.PartitionConsumer, error) {
	logger.Info("Kafka consumer is running on:", zap.String("host:port", a.diProvider.config.KafkaConfig.ConnectString()))
	partConsumer, err := a.kafkaConsumer.ConsumePartition("chat123", 0, sarama.OffsetNewest)
	if err != nil {
		logger.Fatal("Failed to consume partition", zap.Any("err", err))
		return nil, err
	}

	return partConsumer, nil
}

func (a *App) runGrpcServer(ctx context.Context) error {
	logger.Info("GRPC server is running on:", zap.String("host:port", a.diProvider.config.GRPCConfig.Address()))

	list, err := net.Listen("tcp", a.diProvider.Config(ctx).GRPCConfig.Address())
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initPrometheus(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	a.prometheus = &http.Server{
		Addr:    a.diProvider.Config(ctx).PrometheusConfig.Address(),
		Handler: mux,
	}

	return prometheusMetrics.Init(ctx)
}

func (a *App) runPrometheus(_ context.Context) error {
	logger.Info("Prometheus server is running on:", zap.String("host:port", a.diProvider.config.PrometheusConfig.Address()))

	err := a.prometheus.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initTrace(_ context.Context) error {
	tracing.Init(logger.Logger(), "chat")

	return nil
}
