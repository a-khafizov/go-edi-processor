package ports

import (
	"context"

	"github.com/go-edi-document-processor/internal/domain"
)

type DocumentRepository interface {
	Create(ctx context.Context, doc *domain.Document) error
	FindByID(ctx context.Context, id string) (*domain.Document, error)
}
