package services

import (
	"context"
	"time"

	"github.com/go-edi-document-processor/internal/core/domain"
	ports "github.com/go-edi-document-processor/internal/core/ports/secondary"
	"github.com/google/uuid"
)

type DocumentService struct {
	documentRepository ports.DocumentRepository
	outboxService      ports.OutboxService
}

func NewDocumentService(documentRepository ports.DocumentRepository, outboxService ports.OutboxService) *DocumentService {
	return &DocumentService{
		documentRepository: documentRepository,
		outboxService:      outboxService,
	}
}

func (s *DocumentService) SendDocument(ctx context.Context, document *domain.Document) (*domain.Document, error) {
	document.DocId = uuid.New().String()
	document.CreatedAt = time.Now()
	document.UpdatedAt = time.Now()
	document.Status = domain.Received

	err := s.outboxService.SaveDocumentWithEvent(ctx, document, "document.send")
	if err != nil {
		return nil, err
	}

	savedDoc := &domain.Document{
		DocId:  document.DocId,
		Status: document.Status,
	}

	return savedDoc, nil
}

func (s *DocumentService) GetDocumentByUUID(ctx context.Context, docId string) (*domain.Document, error) {
	return s.documentRepository.Get(ctx, docId)
}

func (s *DocumentService) ReceiveDocument(ctx context.Context) (*domain.Document, error) {
	doc, err := s.documentRepository.GetOldest(ctx)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, nil
	}

	doc.Status = domain.Successful
	doc.UpdatedAt = time.Now()

	err = s.outboxService.SaveDocumentWithEvent(ctx, doc, "document.status.updated")
	if err != nil {
		return nil, err
	}

	return doc, nil
}
