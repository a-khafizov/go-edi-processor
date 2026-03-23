package adapters

import (
	"context"

	"github.com/go-edi-document-processor/internal/core/domain"
)

type DocumentRepository struct {
}

func NewDocumentRepository() *DocumentRepository {
	return &DocumentRepository{}
}

func (r *DocumentRepository) Save(ctx context.Context, doc *domain.Document) error {
	return nil
}

func (r *DocumentRepository) Get(ctx context.Context, id string) (*domain.Document, error) {
	return nil, nil
}
