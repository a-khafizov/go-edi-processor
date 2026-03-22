package domain

import "time"

type Document struct {
	ID         string
	Type       DocumentType
	Content    []byte
	SenderID   string
	ReceiverID string
	Status     DocumentStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type DocumentType string

const (
	DocumentTypeXML  DocumentType = "XML"
	DocumentTypePDF  DocumentType = "PDF"
	DocumentTypeJSON DocumentType = "JSON"
)

type DocumentStatus string

const (
	DocumentStatusPending   DocumentStatus = "PENDING"
	DocumentStatusReceived  DocumentStatus = "RECEIVED"
	DocumentStatusProcessed DocumentStatus = "PROCESSED"
	DocumentStatusFailed    DocumentStatus = "FAILED"
)
