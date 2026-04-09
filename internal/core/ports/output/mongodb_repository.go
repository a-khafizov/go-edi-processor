package ports

import (
	"context"

	"github.com/go-edi-document-processor/internal/core/domain"
)

type MongoDocumentRepository interface {
	Insert(ctx context.Context, doc *domain.Document) error
}
