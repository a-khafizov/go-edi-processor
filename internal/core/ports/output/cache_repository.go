package ports

import (
	"context"
	"time"

	"github.com/go-edi-document-processor/internal/core/domain"
)

type CacheRepository interface {
	Set(ctx context.Context, key string, doc *domain.Document, ttl time.Duration) error
	Get(ctx context.Context, key string) (*domain.Document, error)
	Delete(ctx context.Context, key string) error
}
