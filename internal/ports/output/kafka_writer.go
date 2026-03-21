package ports

import "context"

type KafkaWriter interface {
	WriteMessage(ctx context.Context, topic string, key, value []byte, headers map[string]string) error
	Close() error
}
