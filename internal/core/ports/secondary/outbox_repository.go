package ports

import (
	"context"

	"github.com/go-edi-document-processor/internal/core/domain"
)

type OutboxRepository interface {
	Save(ctx context.Context, doc *domain.OutboxMessage) error
}
