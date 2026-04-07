package adapters

import (
	"context"
	"errors"
	"time"

	"github.com/go-edi-document-processor/internal/core/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (m *documentModel) toDomain() *domain.Document {
	return &domain.Document{
		DocId:      m.DocId,
		Type:       domain.DocumentType(m.Type),
		Content:    m.Content,
		SenderID:   m.SenderID,
		ReceiverID: m.ReceiverID,
		Status:     domain.DocumentStatus(m.Status),
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
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

func (r *MongoDocumentRepository) FindByID(ctx context.Context, id string) (*domain.Document, error) {
	var model documentModel

	err := r.collection.FindOne(ctx, bson.M{"doc_id": id}).Decode(&model)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return model.toDomain(), nil
}

func (r *MongoDocumentRepository) FindAll(ctx context.Context, limit, skip int64) ([]*domain.Document, error) {
	opts := options.Find().
		SetLimit(limit).
		SetSkip(skip).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []*domain.Document

	for cursor.Next(ctx) {
		var model documentModel
		if err := cursor.Decode(&model); err != nil {
			return nil, err
		}
		docs = append(docs, model.toDomain())
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return docs, nil
}

func (r *MongoDocumentRepository) Update(ctx context.Context, doc *domain.Document) error {
	model := fromDomain(doc)

	filter := bson.M{"doc_id": doc.DocId}
	update := bson.M{"$set": model}

	_, err := r.collection.UpdateOne(ctx, filter, update)

	return err
}

func (r *MongoDocumentRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"doc_id": id})

	return err
}

func (r *MongoDocumentRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}

func (r *MongoDocumentRepository) Ping(ctx context.Context) error {
	return r.client.Ping(ctx, nil)
}
