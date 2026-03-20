package domain

import "time"

type OutboxMessage struct {
	ID          int64
	Topic       string
	Message     []byte
	Key         []byte
	Headers     map[string]string
	CreatedAt   time.Time
	ProcessedAt *time.Time
	Delay       time.Duration
}
