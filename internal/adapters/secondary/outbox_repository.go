package adapters

import (
	"context"

	"github.com/go-edi-document-processor/internal/core/domain"
)

type OutboxRepository struct {
}

func NewOutboxRepository() *OutboxRepository {
	return &OutboxRepository{}
}

func (r *OutboxRepository) Save(ctx context.Context, message *domain.OutboxMessage) error {
	return nil
}
