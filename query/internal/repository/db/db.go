package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

type Repository struct {
	conn *mongo.Client
}

func New(mongoAddr string) *Repository {
	ctx := context.Background()
	clientOpts := options.Client().ApplyURI(mongoAddr)
	mongoCli, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal(err)
	}
	if err = mongoCli.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	return &Repository{conn: mongoCli}
}

func (r *Repository) Close() {
	err := r.conn.Disconnect(context.Background())
	if err != nil {
		log.Println("Failed to close connection to MongoDB: ", err)
	}
}

func (r *Repository) UserPage(ctx context.Context, userID uint64) ([]byte, error) {
	return nil, nil
}
