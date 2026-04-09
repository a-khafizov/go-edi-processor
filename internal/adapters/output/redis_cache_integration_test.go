//go:build integration

package adapters

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-edi-document-processor/internal/core/domain"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisCache_Integration(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("could not start miniredis: %v", err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer client.Close()

	cache := NewRedisCache(client)
	ctx := context.Background()

	t.Run("Set and Get document", func(t *testing.T) {
		doc := &domain.Document{
			DocId:      "test-doc-1",
			Type:       domain.XML,
			Content:    []byte("<xml>test</xml>"),
			SenderID:   "sender1",
			ReceiverID: "receiver1",
			Status:     domain.Pending,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		err := cache.Set(ctx, doc.DocId, doc, 5*time.Minute)
		assert.NoError(t, err)

		retrieved, err := cache.Get(ctx, doc.DocId)
		assert.NoError(t, err)
		assert.Equal(t, doc.DocId, retrieved.DocId)
		assert.Equal(t, doc.Type, retrieved.Type)
		assert.Equal(t, doc.Status, retrieved.Status)
		assert.Equal(t, doc.SenderID, retrieved.SenderID)
		assert.Equal(t, doc.ReceiverID, retrieved.ReceiverID)
		assert.Equal(t, doc.Content, retrieved.Content)
	})

	t.Run("Get non-existent key returns nil", func(t *testing.T) {
		doc, err := cache.Get(ctx, "non-existent")
		assert.NoError(t, err)
		assert.Nil(t, doc)
	})

	t.Run("Delete key", func(t *testing.T) {
		doc := &domain.Document{
			DocId:  "to-delete",
			Type:   domain.JSON,
			Status: domain.Received,
		}
		err := cache.Set(ctx, doc.DocId, doc, time.Minute)
		assert.NoError(t, err)

		retrieved, err := cache.Get(ctx, doc.DocId)
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)

		err = cache.Delete(ctx, doc.DocId)
		assert.NoError(t, err)

		retrieved, err = cache.Get(ctx, doc.DocId)
		assert.NoError(t, err)
		assert.Nil(t, retrieved)
	})

	t.Run("Set with zero TTL", func(t *testing.T) {
		doc := &domain.Document{DocId: "zero-ttl", Type: domain.PDF}
		err := cache.Set(ctx, doc.DocId, doc, 0)
		assert.NoError(t, err)

		retrieved, err := cache.Get(ctx, doc.DocId)
		assert.NoError(t, err)
		assert.Equal(t, doc.DocId, retrieved.DocId)
	})

	t.Run("JSON marshalling error simulation", func(t *testing.T) {
	})
}
