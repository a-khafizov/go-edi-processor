package ports

import (
	"context"

	"github.com/go-edi-document-processor/internal/domain"
)

// OutboxRepository интерфейс репозитория для Outbox сообщений
type OutboxRepository interface {
	Create(ctx context.Context, msg *domain.OutboxMessage) error
	GetUnprocessed(ctx context.Context, limit int) ([]*domain.OutboxMessage, error)
	MarkAsProcessed(ctx context.Context, id int64) error
	DeleteProcessed(ctx context.Context, olderThanDays int) error
}
