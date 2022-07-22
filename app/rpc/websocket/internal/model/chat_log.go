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

type ChatLog struct {
	ServerMsgID    primitive.ObjectID `bson:"_id,omitempty"`
	ConversationID primitive.ObjectID `bson:"conversation_id"`
	ClientMsgID    string             `bson:"client_msg_id"`
	ClientTime     int64              `bson:"client_time"`
	ServerTime     int64              `bson:"server_time"`
	SenderID       string             `bson:"sender_id"`
	Seq            uint32             `bson:"seq"`
	IsRecalled     bool               `bson:"is_recalled"`
	RecallTips     []byte             `bson:"recall_tips"`
	DeleteUserIds  []string           `bson:"delete_user_ids"`
	Data           ChatLogData        `bson:"data"`
	OfflinePush    OfflinePush        `bson:"offline_push"`
	MsgOptions     MsgOptions         `bson:"msg_options"`
}

type (
	ChatLogData struct {
		ContentType int32  `bson:"content_type"`
		Content     []byte `bson:"content"`
	}
	OfflinePush struct {
		Title         string `bson:"title"`
		Desc          string `bson:"desc"`
		Ex            string `bson:"ex"`
		IOSPushSound  string `bson:"IOSPushSound"`
		IOSBadgeCount bool   `bson:"IOSBadgeCount"`
	}
	MsgOptions struct {
		Storage bool `bson:"storage"`
		Unread  bool `bson:"unread"`
	}
)

func NewOfflinePush(pb *pb.OfflinePush) OfflinePush {
	if pb == nil {
		return OfflinePush{}
	}
	return OfflinePush{
		Title:         pb.Title,
		Desc:          pb.Desc,
		Ex:            pb.Ex,
		IOSPushSound:  pb.IOSPushSound,
		IOSBadgeCount: pb.IOSBadgeCount,
	}
}

func NewMsgOptions(pb *pb.MsgOptions) MsgOptions {
	if pb == nil {
		return MsgOptions{}
	}
	return MsgOptions{
		Storage: pb.Storage,
		Unread:  pb.Unread,
	}
}

func (m *ChatLog) CollectionName() string {
	return "chat_log"
}

type IndexView struct {
	V    int                    `json:"v"`
	Key  map[string]interface{} `json:"key"`
	Name string                 `json:"name"`
	Ns   string                 `json:"ns"`
}

func (m *ChatLog) Indexes(c *mongo.Collection) error {
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
	keysDoc := bsonx.Doc{}
	createMany, err := c.Indexes().CreateMany(
		context.Background(),
		[]mongo.IndexModel{{
			Keys: keysDoc.Append("conversation_id", bsonx.Int64(1)).Append("seq", bsonx.Int64(1)),
			Options: &options.IndexOptions{
				Unique: &yes,
			},
		}, {
			Keys: bson.M{
				"conversation_id": "hashed",
			},
		}, {
			Keys: bson.M{
				"seq": -1,
			},
		}, {
			Keys: bson.M{
				"client_msg_id": "hashed",
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
