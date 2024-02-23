package serviceHistory

import (
	"context"
	"github.com/PerfilievAlexandr/chat-server/internal/domain"
	"github.com/PerfilievAlexandr/chat-server/internal/repository"
	"github.com/PerfilievAlexandr/chat-server/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type historyService struct {
	historyRepository repository.HistoryRepository
}

func NewHistoryService(_ context.Context, historyRepository repository.HistoryRepository) service.HistoryService {
	return &historyService{
		historyRepository: historyRepository,
	}
}

func (h *historyService) SaveHistory(ctx context.Context, message domain.Message) error {
	err := h.historyRepository.SaveHistory(ctx, message)
	if err != nil {
		return status.Errorf(codes.Internal, "error save history")
	}

	return nil
}
