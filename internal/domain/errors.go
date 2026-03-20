package domain

import "errors"

var (
	ErrInvalidDocumentID        = errors.New("invalid document ID")
	ErrInvalidDocumentType      = errors.New("invalid document type")
	ErrEmptyContent             = errors.New("document content cannot be empty")
	ErrInvalidSender            = errors.New("sender ID cannot be empty")
	ErrInvalidReceiver          = errors.New("receiver ID cannot be empty")
	ErrDocumentNotFound         = errors.New("document not found")
	ErrDocumentAlreadyProcessed = errors.New("document already processed")
	ErrOutboxMessageNotFound    = errors.New("outbox message not found")
)
