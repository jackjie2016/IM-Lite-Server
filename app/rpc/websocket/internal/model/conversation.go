package model

import (
	"context"
	"github.com/globalsign/mgo/bson"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var yes = true

type Conversation struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Seq       uint64             `bson:"seq"`
	Members   []string           `bson:"members"`
	DeletedAt int64              `bson:"deleted_at"`
}

func (m *Conversation) Indexes(c *mongo.Collection) error {
	// 查询索引
	indexes, err := c.Indexes().List(context.TODO())
	if err != nil {
		logx.Errorf("get chat_log indexes error: %v", err)
		return err
	}
	var indexViews []IndexView
	err = indexes.All(context.Background(), &indexViews)
	if err != nil {
		logx.Errorf("get chat_log indexes error: %v", err)
		return err
	}
	if len(indexViews) > 1 {
		return nil
	}
	createMany, err := c.Indexes().CreateMany(
		context.Background(),
		[]mongo.IndexModel{{
			Keys: bson.M{
				"deleted_at": 1,
			},
		}},
	)
	if err != nil {
		logx.Errorf("create index error: %v", err)
		return nil
	} else {
		logx.Infof("create index: %+v", createMany)
	}
	return nil
}

func (m *Conversation) CollectionName() string {
	return "conversation"
}
