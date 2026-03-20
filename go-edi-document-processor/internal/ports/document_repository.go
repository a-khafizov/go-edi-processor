package ports

import (
	"context"

	"github.com/go-edi-document-processor/internal/domain"
)

type DocumentRepository interface {
	Create(ctx context.Context, doc *domain.Document) error
	FindByID(ctx context.Context, id string) (*domain.Document, error)
	Update(ctx context.Context, doc *domain.Document) error
	Delete(ctx context.Context, id string) error
	ListByStatus(ctx context.Context, status domain.DocumentStatus, limit, offset int) ([]*domain.Document, error)
}
