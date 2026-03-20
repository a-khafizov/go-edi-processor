package ports

import (
	"context"

	"github.com/go-edi-document-processor/internal/domain"
)

// DocumentService интерфейс сервиса обработки документов
type DocumentService interface {
	SubmitDocument(ctx context.Context, doc *domain.Document) error
	GetDocument(ctx context.Context, id string) (*domain.Document, error)
	ProcessDocument(ctx context.Context, id string) error
	ListDocuments(ctx context.Context, status domain.DocumentStatus, limit, offset int) ([]*domain.Document, error)
}
