package db

import (
	"context"
	"time"

	"github.com/p4u/padelfriends/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB holds the Mongo client and the main database reference
type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// Connect initializes a MongoDB client and selects a database.
// Database name can be hardcoded or configurable. We'll hardcode "padelfriends" for now.
func Connect(cfg *config.Config) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	// Ping to ensure connection is established
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := client.Database("padelfriends")

	return &MongoDB{
		Client:   client,
		Database: db,
	}, nil
}
