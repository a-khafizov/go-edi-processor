package ports

import (
	"context"

	"github.com/go-edi-document-processor/internal/core/domain"
)

type OutboxService interface {
	SaveDocumentWithEvent(ctx context.Context, doc *domain.Document, eventType string) error
}
