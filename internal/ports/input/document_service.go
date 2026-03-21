package ports

import (
	"context"

	"github.com/go-edi-document-processor/internal/domain"
)

type DocumentService interface {
	GetDocumentByUUID(ctx context.Context, uuid string) (*domain.Document, error)
	SendDocument(ctx context.Context, doc *domain.Document) error
}
