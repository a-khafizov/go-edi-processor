package deps

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func InitMongoDB(cfg *Config) (*mongo.Client, error) {
	if !cfg.MongoDBEnabled {
		log.Println("MongoDB is disabled, skipping initialization")
		return nil, nil
	}

	var uri string
	var credential options.Credential

	authSource := cfg.MongoDBAuthSource

	if cfg.MongoDBRootUser != "" && cfg.MongoDBRootPassword != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s",
			cfg.MongoDBRootUser,
			cfg.MongoDBRootPassword,
			cfg.MongoDBHost,
			cfg.MongoDBPort,
		)
		credential = options.Credential{
			Username:   cfg.MongoDBRootUser,
			Password:   cfg.MongoDBRootPassword,
			AuthSource: authSource,
		}
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s",
			cfg.MongoDBHost,
			cfg.MongoDBPort,
		)
	}

	clientOptions := options.Client().ApplyURI(uri)
	if credential.Username != "" {
		clientOptions.SetAuth(credential)
	}
	clientOptions.SetConnectTimeout(10 * time.Second)
	clientOptions.SetServerSelectionTimeout(10 * time.Second)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		client.Disconnect(ctx)
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("MongoDB connected successfully")
	return client, nil
}
