package mongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	ConnectionString string `required:"true" split_words:"true"`
	Database         string `required:"true" split_words:"true"`
}

type Mongo struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func New(config Config) (Mongo, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.ConnectionString))
	if err != nil {
		return Mongo{}, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return Mongo{}, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return Mongo{}, err
	}
	db := client.Database(config.Database)
	return Mongo{client, db}, nil
}

func (m Mongo) Disconnect() {
	err := m.Client.Disconnect(context.Background())
	if err != nil {
		log.Println("error in disconnecting ")
	}
}
