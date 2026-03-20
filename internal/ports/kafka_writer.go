package ports

import "context"

// KafkaWriter интерфейс для отправки сообщений в Kafka
type KafkaWriter interface {
	WriteMessage(ctx context.Context, topic string, key, value []byte, headers map[string]string) error
	Close() error
}
