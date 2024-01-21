package app

import (
	"context"
	api "github.com/PerfilievAlexandr/chat-server/internal/api/grpc/message"
	"github.com/PerfilievAlexandr/chat-server/internal/config"
	"github.com/PerfilievAlexandr/chat-server/internal/repository"
	repo "github.com/PerfilievAlexandr/chat-server/internal/repository/message"
	"github.com/PerfilievAlexandr/chat-server/internal/service"
	serviceImpl "github.com/PerfilievAlexandr/chat-server/internal/service/message"
	"github.com/PerfilievAlexandr/platform_common/pkg/closer"
	"github.com/PerfilievAlexandr/platform_common/pkg/db"
	"github.com/PerfilievAlexandr/platform_common/pkg/db/pg"
	"github.com/PerfilievAlexandr/platform_common/pkg/db/transaction"
	"log"
)

type diProvider struct {
	config         *config.Config
	db             db.Client
	txManager      db.TxManager
	messageRepo    repository.MessageRepository
	messageService service.MessageService
	messageServer  *api.Server
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
		p.messageRepo = repo.NewMessageRepo(ctx, p.Db(ctx))
	}

	return p.messageRepo
}

func (p *diProvider) MessageService(ctx context.Context) service.MessageService {
	if p.messageService == nil {
		p.messageService = serviceImpl.NewMessageService(ctx, p.MessageRepository(ctx), p.TxManager(ctx))
	}

	return p.messageRepo
}

func (p *diProvider) MessageServer(ctx context.Context) *api.Server {
	if p.messageServer == nil {
		p.messageServer = api.NewServer(ctx, p.MessageService(ctx))
	}

	return p.messageServer
}
