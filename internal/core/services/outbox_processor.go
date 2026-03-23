package services

import (
	"context"

	"github.com/go-edi-document-processor/internal/core/domain"
	ports "github.com/go-edi-document-processor/internal/core/ports/secondary"
)

type OutboxService struct {
	outboxRepository ports.OutboxRepository
}

func NewOutboxService(outboxRepository ports.OutboxRepository) *OutboxService {
	return &OutboxService{
		outboxRepository: outboxRepository,
	}
}

func (s *OutboxService) Save(ctx context.Context, message *domain.OutboxMessage) error {
	return s.outboxRepository.Save(ctx, message)
}
