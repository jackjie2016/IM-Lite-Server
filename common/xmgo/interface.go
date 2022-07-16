package xmgo

import "go.mongodb.org/mongo-driver/mongo"

type ICollection interface {
	CollectionName() string
	Indexes(c *mongo.Collection) error
}
