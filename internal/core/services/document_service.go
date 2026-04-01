package services

import (
	"context"
	"fmt"
	"time"

	"github.com/go-edi-document-processor/internal/core/domain"
	ports "github.com/go-edi-document-processor/internal/core/ports/secondary"
	"github.com/google/uuid"
)

type DocumentService struct {
	documentRepository ports.DocumentRepository
	outboxService      ports.OutboxService
	cacheRepository    ports.CacheRepository
	mongoRepository    ports.MongoDocumentRepository // опционально, может быть nil
}

func NewDocumentService(documentRepository ports.DocumentRepository, outboxService ports.OutboxService, cacheRepository ports.CacheRepository, mongoRepository ports.MongoDocumentRepository) *DocumentService {
	return &DocumentService{
		documentRepository: documentRepository,
		outboxService:      outboxService,
		cacheRepository:    cacheRepository,
		mongoRepository:    mongoRepository,
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

	if err := s.cacheRepository.Set(ctx, document.DocId, document, 5*time.Minute); err != nil {
		fmt.Printf("Warning: failed to cache document %s: %v\n", document.DocId, err)
	}

	if s.mongoRepository != nil {
		go func() {
			mongoCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			if err := s.mongoRepository.Insert(mongoCtx, document); err != nil {
				fmt.Printf("Warning: failed to save document %s to MongoDB: %v\n", document.DocId, err)
			}
		}()
	}

	savedDoc := &domain.Document{
		DocId:  document.DocId,
		Status: document.Status,
	}

	return savedDoc, nil
}

func (s *DocumentService) GetDocumentByUUID(ctx context.Context, docId string) (*domain.Document, error) {
	cachedDoc, err := s.cacheRepository.Get(ctx, docId)
	if err != nil {
		return nil, fmt.Errorf("failed to get from cache: %w", err)
	}
	if cachedDoc != nil {
		return cachedDoc, nil
	}

	doc, err := s.documentRepository.Get(ctx, docId)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, nil
	}

	if err := s.cacheRepository.Set(ctx, docId, doc, 5*time.Minute); err != nil {
		fmt.Printf("Warning: failed to cache document %s: %v\n", docId, err)
	}

	return doc, nil
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

	if err := s.cacheRepository.Delete(ctx, doc.DocId); err != nil {
		fmt.Printf("Warning: failed to invalidate cache for document %s: %v\n", doc.DocId, err)
	}

	return doc, nil
}
