package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-edi-document-processor/internal/deps"
	"github.com/oagudo/outbox"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kmsg"
	"github.com/twmb/franz-go/pkg/sasl/plain"
)

type KafkaPublisher struct {
	client *kgo.Client
	topic  string
}

func NewKafkaPublisher(cfg *deps.Config) (*KafkaPublisher, error) {
	if cfg.KafkaTopic == "" {
		return nil, fmt.Errorf("Kafka topic is not configured")
	}

	brokers := strings.Split(cfg.KafkaBrokers, ",")
	if len(brokers) == 0 {
		return nil, fmt.Errorf("no Kafka brokers configured")
	}

	opts := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
	}

	securityProtocol := strings.ToUpper(cfg.KafkaSecurityProtocol)
	if securityProtocol == "SASL_PLAINTEXT" || securityProtocol == "PLAINTEXT" {
		opts = append(opts, kgo.DialTLSConfig(nil))
	}

	if cfg.KafkaUsername != "" && cfg.KafkaPassword != "" {
		mechanism := plain.Auth{
			User: cfg.KafkaUsername,
			Pass: cfg.KafkaPassword,
		}.AsMechanism()
		opts = append(opts, kgo.SASL(mechanism))
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka client: %w", err)
	}

	if err := ensureTopic(client, cfg.KafkaTopic); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to ensure Kafka topic: %w", err)
	}

	return &KafkaPublisher{
		client: client,
		topic:  cfg.KafkaTopic,
	}, nil
}

func ensureTopic(client *kgo.Client, topic string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	metaReq := kmsg.NewMetadataRequest()
	metaReqTopic := kmsg.NewMetadataRequestTopic()
	metaReqTopic.Topic = kmsg.StringPtr(topic)
	metaReq.Topics = []kmsg.MetadataRequestTopic{metaReqTopic}

	metaResp, err := client.Request(ctx, &metaReq)
	if err != nil {
		return fmt.Errorf("failed to get metadata: %w", err)
	}
	metadata := metaResp.(*kmsg.MetadataResponse)

	for _, t := range metadata.Topics {
		if t.Topic != nil && *t.Topic == topic && t.ErrorCode == 0 {
			return nil
		}
	}

	createReq := kmsg.NewCreateTopicsRequest()
	createReqTopic := kmsg.NewCreateTopicsRequestTopic()
	createReqTopic.Topic = topic
	createReqTopic.NumPartitions = 1
	createReqTopic.ReplicationFactor = 1
	createReq.Topics = []kmsg.CreateTopicsRequestTopic{createReqTopic}

	createResp, err := client.Request(ctx, &createReq)
	if err != nil {
		return fmt.Errorf("failed to send create topics request: %w", err)
	}
	createResponse := createResp.(*kmsg.CreateTopicsResponse)

	for _, t := range createResponse.Topics {
		if t.Topic == topic && t.ErrorCode != 0 {
			return fmt.Errorf("failed to create topic %s: error code %d", topic, t.ErrorCode)
		}
	}

	return nil
}

func (p *KafkaPublisher) Publish(ctx context.Context, msg *outbox.Message) error {
	var payload map[string]interface{}
	var key []byte
	if err := json.Unmarshal(msg.Payload, &payload); err == nil {
		if docID, ok := payload["doc_id"].(string); ok {
			key = []byte(docID)
		}
	}
	record := &kgo.Record{
		Topic: p.topic,
		Key:   key,
		Value: msg.Payload,
	}
	err := p.client.ProduceSync(ctx, record).FirstErr()
	if err != nil {
		return fmt.Errorf("failed to publish message to Kafka: %w", err)
	}
	return nil
}

func (p *KafkaPublisher) Close() {
	p.client.Close()
}
