package adapters

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-edi-document-processor/internal/core/domain"
	"github.com/oagudo/outbox"
)

type DocumentRepository struct {
	db *sql.DB
}

func NewDocumentRepository(db *sql.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

func (r *DocumentRepository) SaveWithTx(ctx context.Context, tx outbox.TxQueryer, doc *domain.Document) error {
	query := `
		insert into documents (doc_id, type, content, sender_id, receiver_id, status, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8)
		on conflict (doc_id) do update set
			type = excluded.type,
			content = excluded.content,
			sender_id = excluded.sender_id,
			receiver_id = excluded.receiver_id,
			status = excluded.status,
			created_at = excluded.created_at,
			updated_at = excluded.updated_at
	`
	_, err := tx.ExecContext(ctx, query,
		doc.DocId,
		string(doc.Type),
		doc.Content,
		doc.SenderID,
		doc.ReceiverID,
		string(doc.Status),
		doc.CreatedAt,
		doc.UpdatedAt,
	)
	return err
}

func (r *DocumentRepository) Get(ctx context.Context, id string) (*domain.Document, error) {
	query := `
		select doc_id, type, content, sender_id, receiver_id, status, created_at, updated_at
		from documents
		where doc_id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	var doc domain.Document
	var docType, status string
	var createdAt, updatedAt time.Time
	err := row.Scan(&doc.DocId, &docType, &doc.Content, &doc.SenderID, &doc.ReceiverID, &status, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	doc.Type = domain.DocumentType(docType)
	doc.Status = domain.DocumentStatus(status)
	doc.CreatedAt = createdAt
	doc.UpdatedAt = updatedAt
	return &doc, nil
}

func (r *DocumentRepository) GetOldest(ctx context.Context) (*domain.Document, error) {
	query := `
		select doc_id, type, content, sender_id, receiver_id, status, created_at, updated_at
		from documents
		where status = $1
		order by created_at asc
		limit 1
	`
	row := r.db.QueryRowContext(ctx, query, domain.Received)
	var doc domain.Document
	var docType, status string
	var createdAt, updatedAt time.Time
	err := row.Scan(&doc.DocId, &docType, &doc.Content, &doc.SenderID, &doc.ReceiverID, &status, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	doc.Type = domain.DocumentType(docType)
	doc.Status = domain.DocumentStatus(status)
	doc.CreatedAt = createdAt
	doc.UpdatedAt = updatedAt
	return &doc, nil
}
