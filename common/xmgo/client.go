package xmgo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Uri string
}

func GetClient(
	cfg Config,
) *mongo.Client {
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.Uri))
	if err != nil {
		panic(err)
	}
	err = mongoClient.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	return mongoClient
}
