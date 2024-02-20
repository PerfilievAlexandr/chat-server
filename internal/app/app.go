package app

import (
	"context"
	"github.com/PerfilievAlexandr/chat-server/internal/api/grpc/chat/interceptor"
	"github.com/PerfilievAlexandr/chat-server/internal/config"
	"github.com/PerfilievAlexandr/chat-server/internal/logger"
	"github.com/PerfilievAlexandr/chat-server/internal/tracing"
	proto "github.com/PerfilievAlexandr/chat-server/pkg/chat_v1"
	"github.com/PerfilievAlexandr/platform_common/pkg/closer"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type App struct {
	diProvider *diProvider
	grpcServer *grpc.Server
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

	return a.runGrpcServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initTrace,
		a.initProvider,
		a.initGrpcServer,
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
			),
		),
	)
	reflection.Register(a.grpcServer)
	proto.RegisterChatV1Server(a.grpcServer, a.diProvider.ChatServer(ctx))

	return nil
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

func (a *App) initTrace(_ context.Context) error {
	tracing.Init("chat")

	return nil
}
