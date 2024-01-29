package app

import (
	"context"
	"github.com/PerfilievAlexandr/chat-server/internal/api/grpc/message/interceptor"
	"github.com/PerfilievAlexandr/chat-server/internal/config"
	proto "github.com/PerfilievAlexandr/chat-server/pkg/chat_v1"
	"github.com/PerfilievAlexandr/platform_common/pkg/closer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
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

	return a.runGRPCServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initProvider,
		a.initGRPCServer,
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

func (a *App) initProvider(_ context.Context) error {
	a.diProvider = newDiProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(interceptor.AccessInterceptor),
	)

	reflection.Register(a.grpcServer)

	proto.RegisterChatV1Server(a.grpcServer, a.diProvider.MessageServer(ctx))

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	log.Printf("GRPC server is running on %s", a.diProvider.Config(ctx).GRPCConfig.Address())

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
