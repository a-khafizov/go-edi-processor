package ports

import (
	"context"

	"github.com/go-edi-document-processor/internal/core/domain"
)

type DocumentRepository interface {
	Save(ctx context.Context, doc *domain.Document) error
	Get(ctx context.Context, id string) (*domain.Document, error)
}
