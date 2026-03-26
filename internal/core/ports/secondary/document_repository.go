package ports

import (
	"context"

	"github.com/go-edi-document-processor/internal/core/domain"
	"github.com/oagudo/outbox"
)

type DocumentRepository interface {
	SaveWithTx(ctx context.Context, tx outbox.TxQueryer, doc *domain.Document) error
	Get(ctx context.Context, id string) (*domain.Document, error)
	GetOldest(ctx context.Context) (*domain.Document, error)
}
