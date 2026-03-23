package services

import (
	"context"

	"github.com/go-edi-document-processor/internal/core/domain"
	ports "github.com/go-edi-document-processor/internal/core/ports/secondary"
)

type DocumentService struct {
	documentRepository ports.DocumentRepository
}

func NewDocumentService(documentRepository ports.DocumentRepository) *DocumentService {
	return &DocumentService{
		documentRepository: documentRepository,
	}
}

func (s *DocumentService) SendDocument(ctx context.Context, document *domain.Document) (*domain.Document, error) {
	return nil, nil
}

func (s *DocumentService) GetDocumentByUUID(ctx context.Context, id string) (*domain.Document, error) {
	return nil, nil
}

func (s *DocumentService) ReceiveDocument(ctx context.Context) (*domain.Document, error) {
	return nil, nil
}
