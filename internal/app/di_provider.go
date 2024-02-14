package app

import (
	"context"
	authProto "github.com/PerfilievAlexandr/auth/pkg/access_v1"
	api "github.com/PerfilievAlexandr/chat-server/internal/api/grpc/chat"
	"github.com/PerfilievAlexandr/chat-server/internal/config"
	authClient "github.com/PerfilievAlexandr/chat-server/internal/integration/auth"
	authClientService "github.com/PerfilievAlexandr/chat-server/internal/integration/auth/impl"
	"github.com/PerfilievAlexandr/chat-server/internal/repository"
	chatRepo "github.com/PerfilievAlexandr/chat-server/internal/repository/chat"
	messageRepo "github.com/PerfilievAlexandr/chat-server/internal/repository/message"
	"github.com/PerfilievAlexandr/chat-server/internal/service"
	chatServiceImpl "github.com/PerfilievAlexandr/chat-server/internal/service/chat"
	"github.com/PerfilievAlexandr/platform_common/pkg/closer"
	"github.com/PerfilievAlexandr/platform_common/pkg/db"
	"github.com/PerfilievAlexandr/platform_common/pkg/db/pg"
	"github.com/PerfilievAlexandr/platform_common/pkg/db/transaction"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type diProvider struct {
	config        *config.Config
	db            db.Client
	txManager     db.TxManager
	messageRepo   repository.MessageRepository
	chatRepo      repository.ChatRepository
	chatService   service.ChatService
	authClient    authClient.AuthServiceClient
	messageServer *api.Server
}

func newDiProvider() *diProvider {
	return &diProvider{}
}

func (p *diProvider) Config(ctx context.Context) *config.Config {
	if p.config == nil {
		cfg, err := config.NewConfig(ctx)
		if err != nil {
			if err != nil {
				log.Fatalf("failed to get pg config: %s", err.Error())
			}
		}

		p.config = cfg
	}

	return p.config
}

func (p *diProvider) Db(ctx context.Context) db.Client {
	if p.db == nil {
		dbPool, err := pg.New(ctx, p.Config(ctx).DbConfig.ConnectString())
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}

		err = dbPool.Ping(ctx)
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}

		closer.Add(func() error {
			dbPool.Close()
			return nil
		})

		p.db = dbPool
	}

	return p.db
}

func (p *diProvider) TxManager(ctx context.Context) db.TxManager {
	if p.txManager == nil {
		p.txManager = transaction.NewTransactionManager(p.Db(ctx))
	}

	return p.txManager
}

func (p *diProvider) MessageRepository(ctx context.Context) repository.MessageRepository {
	if p.messageRepo == nil {
		p.messageRepo = messageRepo.NewMessageRepo(ctx, p.Db(ctx))
	}

	return p.messageRepo
}

func (p *diProvider) ChatRepository(ctx context.Context) repository.ChatRepository {
	if p.chatRepo == nil {
		p.chatRepo = chatRepo.NewChatRepo(ctx, p.Db(ctx))
	}

	return p.chatRepo
}

func (p *diProvider) ChatService(ctx context.Context) service.ChatService {
	if p.chatService == nil {
		p.chatService = chatServiceImpl.NewChatService(
			ctx,
			p.MessageRepository(ctx),
			p.ChatRepository(ctx),
			p.TxManager(ctx),
		)
	}

	return p.chatService
}

func (p *diProvider) ChatServer(ctx context.Context) *api.Server {
	if p.messageServer == nil {
		p.messageServer = api.NewServer(ctx, p.ChatService(ctx))
	}

	return p.messageServer
}

func (p *diProvider) AuthClient(ctx context.Context) authClient.AuthServiceClient {
	if p.authClient == nil {
		conn, err := grpc.Dial(
			p.Config(ctx).AuthClientConfig.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
		)
		if err != nil {
			log.Fatalf("failed to dial GRPC client: %v", err)
		}

		p.authClient = authClientService.New(authProto.NewAccessV1Client(conn))
	}

	return p.authClient
}
