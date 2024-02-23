package app

import (
	"context"
	authProto "github.com/PerfilievAlexandr/auth/pkg/access_v1"
	api "github.com/PerfilievAlexandr/chat-server/internal/api/grpc/chat"
	"github.com/PerfilievAlexandr/chat-server/internal/config"
	authClient "github.com/PerfilievAlexandr/chat-server/internal/integration/auth"
	authClientService "github.com/PerfilievAlexandr/chat-server/internal/integration/auth/impl"
	"github.com/PerfilievAlexandr/chat-server/internal/logger"
	"github.com/PerfilievAlexandr/chat-server/internal/repository"
	chatRepo "github.com/PerfilievAlexandr/chat-server/internal/repository/chat"
	historyRepository "github.com/PerfilievAlexandr/chat-server/internal/repository/history"
	messageRepo "github.com/PerfilievAlexandr/chat-server/internal/repository/message"
	"github.com/PerfilievAlexandr/chat-server/internal/service"
	chatServiceImpl "github.com/PerfilievAlexandr/chat-server/internal/service/chat"
	checkRoleServiceImpl "github.com/PerfilievAlexandr/chat-server/internal/service/check_role"
	serviceHistory "github.com/PerfilievAlexandr/chat-server/internal/service/history"
	serviceMessage "github.com/PerfilievAlexandr/chat-server/internal/service/message"
	"github.com/PerfilievAlexandr/platform_common/pkg/closer"
	"github.com/PerfilievAlexandr/platform_common/pkg/db"
	"github.com/PerfilievAlexandr/platform_common/pkg/db/pg"
	"github.com/PerfilievAlexandr/platform_common/pkg/db/transaction"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type diProvider struct {
	config           *config.Config
	db               db.Client
	txManager        db.TxManager
	messageRepo      repository.MessageRepository
	historyRepo      repository.HistoryRepository
	chatRepo         repository.ChatRepository
	chatService      service.ChatService
	messageService   service.MessageService
	historyService   service.HistoryService
	checkRoleService service.CheckRoleService
	authClient       authClient.AuthServiceClient
	messageServer    *api.Server
}

func newDiProvider() *diProvider {
	return &diProvider{}
}

func (p *diProvider) Config(ctx context.Context) *config.Config {
	if p.config == nil {
		cfg, err := config.NewConfig(ctx)
		if err != nil {
			if err != nil {
				logger.Fatal("failed to get pg config", zap.Any("err", err))
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
			logger.Fatal("failed to connect to database", zap.Any("err", err))
		}

		err = dbPool.Ping(ctx)
		if err != nil {
			logger.Fatal("failed to ping database", zap.Any("err", err))
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

func (p *diProvider) HistoryRepository(ctx context.Context) repository.HistoryRepository {
	if p.historyRepo == nil {
		p.historyRepo = historyRepository.NewHistoryRepo(ctx, p.Db(ctx))
	}

	return p.historyRepo
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
			p.MessageService(ctx),
			p.ChatRepository(ctx),
		)
	}

	return p.chatService
}

func (p *diProvider) MessageService(ctx context.Context) service.MessageService {
	if p.messageService == nil {
		p.messageService = serviceMessage.NewMessageService(
			ctx,
			p.MessageRepository(ctx),
			p.HistoryService(ctx),
			p.TxManager(ctx),
		)
	}

	return p.messageService
}

func (p *diProvider) HistoryService(ctx context.Context) service.HistoryService {
	if p.historyService == nil {
		p.historyService = serviceHistory.NewHistoryService(
			ctx,
			p.HistoryRepository(ctx),
		)
	}

	return p.historyService
}

func (p *diProvider) CheckRoleService(ctx context.Context) service.CheckRoleService {
	if p.checkRoleService == nil {
		p.checkRoleService = checkRoleServiceImpl.NewCheckRoleService(ctx)
	}

	return p.checkRoleService
}

func (p *diProvider) ChatServer(ctx context.Context) *api.Server {
	if p.messageServer == nil {
		p.messageServer = api.NewServer(
			ctx,
			p.ChatService(ctx),
			p.CheckRoleService(ctx),
			p.AuthClient(ctx),
		)
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
			logger.Fatal("failed to dial GRPC client:", zap.Any("err", err))
		}

		p.authClient = authClientService.New(authProto.NewAccessV1Client(conn))
	}

	return p.authClient
}
