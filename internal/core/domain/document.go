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
	Pending    DocumentStatus = "PENDING"
	Received   DocumentStatus = "RECEIVED"
	Processed  DocumentStatus = "PROCESSED"
	Failed     DocumentStatus = "FAILED"
	Successful DocumentStatus = "SUCCESSFUL"
)
