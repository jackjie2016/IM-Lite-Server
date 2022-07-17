package model

import (
	"context"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type MidUserConversation struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	// ConversationID 分片健
	ConversationID primitive.ObjectID  `bson:"conversation_id"`
	BelongUser     string              `bson:"belong_user"`
	Remark         string              `bson:"remark"`
	Avatar         string              `bson:"avatar"`
	Unread         uint32              `bson:"unread"`
	RecvOpt        pb.MsgRecvOpt       `bson:"recv_opt"`
	DeletedAt      int64               `bson:"deleted_at"`
	Top            bool                `bson:"top"`
	Timestamp      primitive.Timestamp `bson:"timestamp"`
	MinSeq         uint32              `bson:"min_seq"`
	MaxSeq         uint32              `bson:"-"`
}

func (m *MidUserConversation) CollectionName() string {
	return "mid_user_conversation"
}

func (m *MidUserConversation) Indexes(c *mongo.Collection) error {
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
	if len(indexViews) > 4 {
		return nil
	}
	keysDoc1 := bsonx.Doc{}
	keysDoc2 := bsonx.Doc{}
	keysDoc3 := bsonx.Doc{}
	createMany, err := c.Indexes().CreateMany(
		context.Background(),
		[]mongo.IndexModel{{
			Keys: keysDoc1.Append("conversation_id", bsonx.Int64(1)).Append("belong_user", bsonx.Int64(1)).Append("deleted_at", bsonx.Int64(1)),
			Options: &options.IndexOptions{
				Unique: &yes,
			},
		}, {
			Keys: bson.M{
				"belong_user": "hashed",
			},
		}, {
			Keys: keysDoc2.Append("top", bsonx.Int64(1)).Append("timestamp", bsonx.Int64(-1)),
		}, {
			Keys: keysDoc3.Append("belong_user", bsonx.Int64(1)).Append("deleted_at", bsonx.Int64(1)),
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
