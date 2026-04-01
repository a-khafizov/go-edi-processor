package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-edi-document-processor/internal/deps"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

type KafkaConsumer struct {
	client *kgo.Client
	topic  string
	logger *zap.Logger
}

func NewKafkaConsumer(cfg *deps.Config, logger *zap.Logger) (*KafkaConsumer, error) {
	if cfg.KafkaTopic == "" {
		return nil, fmt.Errorf("Kafka topic is not configured")
	}
	if cfg.KafkaGroupID == "" {
		return nil, fmt.Errorf("Kafka group ID is not configured")
	}

	consumerOpts := []kgo.Opt{
		kgo.ConsumeTopics(cfg.KafkaTopic),
		kgo.ConsumerGroup(cfg.KafkaGroupID),
		kgo.AutoCommitInterval(5 * time.Second),
	}

	opts, err := deps.KafkaClientOptions(cfg, consumerOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka client options: %w", err)
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka client: %w", err)
	}

	return &KafkaConsumer{
		client: client,
		topic:  cfg.KafkaTopic,
		logger: logger,
	}, nil
}

func (c *KafkaConsumer) Start(ctx context.Context) {
	c.logger.Info("Starting Kafka consumer", zap.String("topic", c.topic))
	go c.consumeLoop(ctx)
}

func (c *KafkaConsumer) consumeLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Kafka consumer loop stopped")
			return
		default:
			fetches := c.client.PollFetches(ctx)
			if fetches.IsClientClosed() {
				c.logger.Warn("Kafka client closed, stopping consumer")
				return
			}
			if errs := fetches.Errors(); len(errs) > 0 {
				for _, fetchErr := range errs {
					c.logger.Error("Fetch error", zap.Error(fetchErr.Err))
				}
				continue
			}

			fetches.EachPartition(func(p kgo.FetchTopicPartition) {
				p.EachRecord(c.processRecord)
			})
		}
	}
}

func (c *KafkaConsumer) processRecord(record *kgo.Record) {
	var payload map[string]interface{}
	if err := json.Unmarshal(record.Value, &payload); err != nil {
		c.logger.Error("Failed to unmarshal Kafka message", zap.Error(err), zap.ByteString("value", record.Value))
		return
	}

	c.logger.Info("Received Kafka message",
		zap.String("topic", record.Topic),
		zap.Int32("partition", record.Partition),
		zap.Int64("offset", record.Offset),
		zap.ByteString("key", record.Key),
		zap.Any("payload", payload),
	)

	// здесь может быть добавлена дополнительная обработка
}

func (c *KafkaConsumer) Close() {
	c.client.Close()
	c.logger.Info("Kafka consumer closed")
}
