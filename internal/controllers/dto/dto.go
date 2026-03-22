package dto

import (
	"github.com/go-edi-document-processor/api/proto"
)

type Document struct {
	DocID      string               `json:"doc_id" validate:"required,min=1"`
	Type       proto.DocumentType   `json:"type" validate:"required,min=1"`
	Content    []byte               `json:"content" validate:"required,min=1"`
	SenderID   string               `json:"sender_id"`
	ReceiverID string               `json:"receiver_id" validate:"required,min=1"`
	Status     proto.DocumentStatus `json:"status"`
}
