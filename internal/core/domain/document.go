package domain

import "time"

type Document struct {
	DocId      string
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
	XML  DocumentType = "XML"
	PDF  DocumentType = "PDF"
	JSON DocumentType = "JSON"
)

type DocumentStatus string

const (
	Pending    DocumentStatus = "DOC_STATUS_PENDING"
	Received   DocumentStatus = "DOC_STATUS_RECEIVED"
	Processed  DocumentStatus = "DOC_STATUS_PROCESSED"
	Failed     DocumentStatus = "DOC_STATUS_FAILED"
	Successful DocumentStatus = "DOC_STATUS_SUCCESSFUL"
)
