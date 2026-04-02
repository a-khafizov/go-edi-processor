package services

import (
	"context"
	"testing"
	"time"

	"github.com/go-edi-document-processor/internal/core/domain"
	"github.com/oagudo/outbox"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOutboxService struct {
	mock.Mock
}

func (m *MockOutboxService) SaveDocumentWithEvent(ctx context.Context, doc *domain.Document, event string) error {
	args := m.Called(ctx, doc, event)
	return args.Error(0)
}

type MockCacheRepository struct {
	mock.Mock
}

func (m *MockCacheRepository) Set(ctx context.Context, key string, doc *domain.Document, ttl time.Duration) error {
	args := m.Called(ctx, key, doc, ttl)
	return args.Error(0)
}

func (m *MockCacheRepository) Get(ctx context.Context, key string) (*domain.Document, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Document), args.Error(1)
}

func (m *MockCacheRepository) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCacheRepository) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockCacheRepository) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockDocumentRepository struct {
	mock.Mock
}

func (m *MockDocumentRepository) SaveWithTx(ctx context.Context, tx outbox.TxQueryer, doc *domain.Document) error {
	args := m.Called(ctx, tx, doc)
	return args.Error(0)
}

func (m *MockDocumentRepository) Get(ctx context.Context, id string) (*domain.Document, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Document), args.Error(1)
}

func (m *MockDocumentRepository) GetOldest(ctx context.Context) (*domain.Document, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Document), args.Error(1)
}

func TestDocumentService_SendDocument_Success(t *testing.T) {
	ctx := context.Background()
	doc := &domain.Document{
		Type:       domain.XML,
		Content:    []byte("<xml>test</xml>"),
		SenderID:   "sender1",
		ReceiverID: "receiver1",
	}

	mockOutbox := &MockOutboxService{}
	mockCache := &MockCacheRepository{}

	mockOutbox.On("SaveDocumentWithEvent", ctx, mock.AnythingOfType("*domain.Document"), "document.send").Return(nil)
	mockCache.On("Set", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("*domain.Document"), 5*time.Minute).Return(nil)

	service := NewDocumentService(nil, mockOutbox, mockCache, nil)
	result, err := service.SendDocument(ctx, doc)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.DocId)
	assert.Equal(t, domain.Received, result.Status)
	mockOutbox.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestDocumentService_GetDocumentByUUID_CacheHit(t *testing.T) {
	ctx := context.Background()
	docId := "test-uuid"
	cachedDoc := &domain.Document{
		DocId:  docId,
		Type:   domain.PDF,
		Status: domain.Pending,
	}

	mockCache := &MockCacheRepository{}
	mockCache.On("Get", ctx, docId).Return(cachedDoc, nil)

	service := NewDocumentService(nil, nil, mockCache, nil)
	doc, err := service.GetDocumentByUUID(ctx, docId)

	assert.NoError(t, err)
	assert.Equal(t, cachedDoc, doc)
	mockCache.AssertExpectations(t)
}

func TestDocumentService_ReceiveDocument_Success(t *testing.T) {
	ctx := context.Background()
	oldDoc := &domain.Document{
		DocId:     "doc1",
		Type:      domain.XML,
		Status:    domain.Pending,
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}

	mockRepo := &MockDocumentRepository{}
	mockOutbox := &MockOutboxService{}
	mockCache := &MockCacheRepository{}

	returnedDoc := &domain.Document{
		DocId:     oldDoc.DocId,
		Type:      oldDoc.Type,
		Status:    oldDoc.Status,
		UpdatedAt: oldDoc.UpdatedAt,
	}
	mockRepo.On("GetOldest", ctx).Return(returnedDoc, nil)
	mockOutbox.On("SaveDocumentWithEvent", ctx, mock.AnythingOfType("*domain.Document"), "document.status.updated").Return(nil)
	mockCache.On("Delete", ctx, "doc1").Return(nil)

	service := NewDocumentService(mockRepo, mockOutbox, mockCache, nil)
	doc, err := service.ReceiveDocument(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, doc)
	assert.Equal(t, domain.Successful, doc.Status)
	assert.True(t, doc.UpdatedAt.After(oldDoc.UpdatedAt))
	mockRepo.AssertExpectations(t)
	mockOutbox.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}
