package adapters

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/go-edi-document-processor/internal/core/domain"
	ports "github.com/go-edi-document-processor/internal/core/ports/output"
	"github.com/oagudo/outbox"
)

type outboxServiceImpl struct {
	dbCtx   *outbox.DBContext
	docRepo ports.DocumentRepository
}

func NewOutboxService(db *sql.DB, dbCtx *outbox.DBContext, docRepo ports.DocumentRepository) (ports.OutboxService, error) {
	return &outboxServiceImpl{dbCtx: dbCtx, docRepo: docRepo}, nil
}

func (o *outboxServiceImpl) SaveDocumentWithEvent(ctx context.Context, doc *domain.Document, eventType string) error {
	writer := outbox.NewWriter(o.dbCtx)

	work := func(ctx context.Context, tx outbox.TxQueryer, msgWriter outbox.MessageWriter) error {
		if err := o.docRepo.SaveWithTx(ctx, tx, doc); err != nil {
			return fmt.Errorf("failed to save document: %w", err)
		}

		payload, _ := json.Marshal(map[string]string{
			"doc_id":      doc.DocId,
			"type":        string(doc.Type),
			"status":      string(doc.Status),
			"receiver_id": doc.ReceiverID,
			"content":     string(doc.Content),
			"event":       eventType,
		})
		msg := outbox.NewMessage(payload)
		if err := msgWriter.Store(ctx, msg); err != nil {
			return fmt.Errorf("failed to store outbox message: %w", err)
		}

		return nil
	}

	return writer.Write(ctx, work)
}
