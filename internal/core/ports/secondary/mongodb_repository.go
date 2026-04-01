package ports

import (
	"context"

	"github.com/go-edi-document-processor/internal/core/domain"
)

type MongoDocumentRepository interface {
	Insert(ctx context.Context, doc *domain.Document) error
	FindByID(ctx context.Context, id string) (*domain.Document, error)
	FindAll(ctx context.Context, limit, skip int64) ([]*domain.Document, error)
	Update(ctx context.Context, doc *domain.Document) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
	Ping(ctx context.Context) error
}
