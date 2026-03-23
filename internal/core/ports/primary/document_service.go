package ports

import (
	"context"

	"github.com/go-edi-document-processor/internal/core/domain"
)

type DocumentService interface {
	SendDocument(ctx context.Context, doc *domain.Document) (*domain.Document, error)
	ReceiveDocument(ctx context.Context) (*domain.Document, error)
	GetDocumentByUUID(ctx context.Context, uuid string) (*domain.Document, error)
}
