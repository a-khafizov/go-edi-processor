package adapters

import (
	"context"
	"time"

	"github.com/oagudo/outbox"
	"go.uber.org/zap"
)

type OutboxReader struct {
	reader *outbox.Reader
	logger *zap.Logger
}

func NewOutboxReader(dbCtx *outbox.DBContext, publisher outbox.MessagePublisher, logger *zap.Logger) *OutboxReader {
	reader := outbox.NewReader(dbCtx, publisher,
		outbox.WithInterval(5*time.Second),
		outbox.WithMaxAttempts(5),
		outbox.WithExponentialDelay(5*time.Second, 10*time.Second),
		outbox.WithReadBatchSize(100),
		outbox.WithPublishTimeout(10*time.Second),
	)
	return &OutboxReader{
		reader: reader,
		logger: logger,
	}
}

func (r *OutboxReader) Start() {
	r.logger.Info("Starting outbox reader")
	go func() {
		r.reader.Start()
	}()

	go r.monitorErrors()
	go r.monitorDiscarded()
}

func (r *OutboxReader) Stop(ctx context.Context) error {
	r.logger.Info("Stopping outbox reader")
	return r.reader.Stop(ctx)
}

func (r *OutboxReader) monitorErrors() {
	for err := range r.reader.Errors() {
		r.logger.Error("Outbox reader error", zap.Error(err))
	}
}

func (r *OutboxReader) monitorDiscarded() {
	for msg := range r.reader.DiscardedMessages() {
		r.logger.Warn("Outbox message discarded after max attempts", zap.Any("message", msg))
	}
}
