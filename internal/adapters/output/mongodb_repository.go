package adapters

import (
	"context"
	"time"

	"github.com/go-edi-document-processor/internal/core/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDocumentRepository struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

func NewMongoDocumentRepository(client *mongo.Client, dbName string, collectionName string) *MongoDocumentRepository {
	db := client.Database(dbName)
	collection := db.Collection(collectionName)

	return &MongoDocumentRepository{
		client:     client,
		database:   db,
		collection: collection,
	}
}

type documentModel struct {
	DocId      string    `bson:"doc_id"`
	Type       string    `bson:"type"`
	Content    []byte    `bson:"content"`
	SenderID   string    `bson:"sender_id"`
	ReceiverID string    `bson:"receiver_id"`
	Status     string    `bson:"status"`
	CreatedAt  time.Time `bson:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at"`
}

func fromDomain(doc *domain.Document) *documentModel {
	return &documentModel{
		DocId:      doc.DocId,
		Type:       string(doc.Type),
		Content:    doc.Content,
		SenderID:   doc.SenderID,
		ReceiverID: doc.ReceiverID,
		Status:     string(doc.Status),
		CreatedAt:  doc.CreatedAt,
		UpdatedAt:  doc.UpdatedAt,
	}
}

func (r *MongoDocumentRepository) Insert(ctx context.Context, doc *domain.Document) error {
	model := fromDomain(doc)
	_, err := r.collection.InsertOne(ctx, model)

	return err
}
