package database

import (
	"context"
	"errors"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/model"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/svc"
	"github.com/IM-Lite/IM-Lite-Server/common/utils"
	"github.com/IM-Lite/IM-Lite-Server/common/xredis/rediskey"
	"github.com/IM-Lite/IM-Lite-Server/common/xtrace"
	"github.com/globalsign/mgo/bson"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
)

type Default struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func (l *Default) HasConversation(conversationIDStr string) (bool, error) {
	conversationID, err := primitive.ObjectIDFromHex(conversationIDStr)
	if err != nil {
		l.Errorf("conversationIDStr to ObjectID error, conversationID: %s, err: %s", conversationIDStr, err)
		return false, err
	}
	ctx, _ := context.WithTimeout(l.ctx, time.Second*time.Duration(l.svcCtx.Config.Mongo.DBTimeout))
	conversation := &model.Conversation{}
	err = l.svcCtx.Collection(conversation).FindOne(ctx, bson.M{"_id": conversationID}).Decode(conversation)
	if err != nil {
		// 是否是找不到
		if err == mongo.ErrNoDocuments {
			l.Infof("this conversation is not exist, conversationID: %s", conversationID)
			return false, nil
		}
		l.Errorf("find conversation error, conversationID: %s, err: %s", conversationID, err)
		return false, err
	}
	return true, nil
}

func (l *Default) InsertMayUpdateOneMessage(chatLog *model.ChatLog) error {
	// 新增seq
	conversationIDStr := chatLog.ConversationID.Hex()
	seq, err := l.svcCtx.Redis().Incr(l.ctx, rediskey.ConversationSeq(conversationIDStr)).Result()
	if err != nil {
		l.Errorf("incr conversation seq error, conversationID: %s, err: %s", conversationIDStr, err)
		return err
	}
	chatLog.Seq = uint32(seq)
	// 新增消息
	ctx, _ := context.WithTimeout(l.ctx, time.Second*time.Duration(l.svcCtx.Config.Mongo.DBTimeout))
	var insertResult *mongo.InsertOneResult
	xtrace.StartFuncSpan(ctx, "mongo:InsertMayUpdateOneMessage", func(ctx context.Context) {
		many, updateErr := l.svcCtx.Collection(&model.ChatLog{}).UpdateMany(
			ctx,
			bson.M{"client_msg_id": chatLog.ClientMsgID},
			bson.M{"$set": bson.M{
				"data": chatLog.Data,
			}})
		if updateErr != nil {
			l.Errorf("update chat log error, client_msg_id: %s, err: %s", chatLog.ClientMsgID, updateErr)
		} else {
			l.Infof("update chat log success, client_msg_id: %s, many: %d", chatLog.ClientMsgID, many.ModifiedCount)
		}
		insertResult, err = l.svcCtx.Collection(&model.ChatLog{}).InsertOne(ctx, chatLog)
	})
	if err != nil {
		l.Errorf("insert chat log error, chatLog: %+v, err: %s", chatLog, err)
		return err
	} else {
		chatLog.ServerMsgID, _ = insertResult.InsertedID.(primitive.ObjectID)
	}
	return nil
}

func (l *Default) CreateConversation(conversation *model.Conversation) (string, error) {
	insertOne, err := l.svcCtx.Collection(conversation).InsertOne(l.ctx, conversation)
	if err != nil {
		l.Errorf("insert conversation error: %s", err.Error())
		return "", err
	}
	conversationID, ok := insertOne.InsertedID.(primitive.ObjectID)
	if !ok {
		l.Errorf("insert conversation error: %s", err.Error())
		return "", errors.New("insert conversation error")
	}
	var documents []interface{}
	for _, member := range conversation.Members {
		documents = append(documents,
			&model.MidUserConversation{
				ConversationID: conversationID,
				BelongUser:     member,
				Timestamp:      primitive.Timestamp{T: uint32(time.Now().Unix())},
				MinSeq:         0,
			})
	}
	_, err = l.svcCtx.Collection(&model.MidUserConversation{}).InsertMany(
		l.ctx,
		documents,
	)
	if err != nil {
		l.Errorf("insert mid user conversation error: %s", err.Error())
		return "", err
	}
	return conversationID.Hex(), nil
}

func (l *Default) DeleteConversationAllMessages(conversationID, userID string) error {
	conversationId, err := primitive.ObjectIDFromHex(conversationID)
	if err != nil {
		return err
	}
	result, err := l.svcCtx.Redis().Get(l.ctx, rediskey.ConversationSeq(conversationID)).Result()
	if err != nil {
		l.Errorf("get conversation seq error: %s", err.Error())
		return err
	}
	minSeq := uint64(utils.InterfaceToInt64(result))
	filter := bsonx.Doc{}
	err = l.svcCtx.Collection(&model.MidUserConversation{}).FindOneAndUpdate(
		l.ctx,
		filter.Append("conversation_id", bsonx.ObjectID(conversationId)).Append("belong_user", bsonx.String(userID)).Append("deleted_at", bsonx.Int64(0)),
		bson.M{
			"$set": bson.M{
				"min_seq": minSeq,
			},
		},
	).Err()
	if err != nil {
		l.Errorf("update conversation min seq error: %s", err.Error())
		return err
	}
	return nil
}

func (l *Default) DeleteConversation(conversationID string) error {
	conversationId, err := primitive.ObjectIDFromHex(conversationID)
	if err != nil {
		return err
	}
	{
		ctx, _ := context.WithTimeout(l.ctx, time.Duration(l.svcCtx.Config.Mongo.DBTimeout)*time.Second)
		_, err = l.svcCtx.Collection(&model.MidUserConversation{}).DeleteMany(
			ctx,
			bson.M{
				"conversation_id": conversationId,
			},
		)
		if err != nil {
			l.Errorf("update conversation min seq error: %s", err.Error())
			return err
		}
		_, err = l.svcCtx.Collection(&model.Conversation{}).DeleteOne(
			ctx,
			bson.M{
				"_id": conversationId,
			},
		)
		if err != nil {
			l.Errorf("update conversation min seq error: %s", err.Error())
			return err
		}
	}
	return nil
}

func (l *Default) ListUserConversations(userID string, pageNo int, pageSize int64) ([]*model.MidUserConversation, error) {
	var (
		conversations []*model.MidUserConversation
		seqKeys       []string
		unreadKey     = rediskey.ConversationUnread(userID)
		seqsMap       = make(map[string]int64)
		unreadMap     = make(map[string]string)
	)
	// mongo中查询会话
	{
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(l.svcCtx.Config.Mongo.DBTimeout)*time.Second)
		cursor, err := l.svcCtx.
			Collection(&model.MidUserConversation{}).Find(
			ctx, bson.M{
				"belong_user": userID,
			},
			&options.FindOptions{
				Limit: &pageSize,
				Sort:  bsonx.Doc{}.Append("top", bsonx.Int32(1)).Append("timestamp", bsonx.Int32(-1)),
			},
		)
		if err != nil {
			l.Errorf("ListUserConversation failed, err: %v", err)
			return nil, err
		}
		if err := cursor.All(ctx, &conversations); err != nil {
			l.Errorf("ListUserConversation failed, err: %v", err)
			return nil, err
		}
		if len(conversations) == 0 {
			return nil, nil
		}
	}
	// redis中查询会话未读数 seq
	{
		for _, conversation := range conversations {
			seqKeys = append(seqKeys, rediskey.ConversationSeq(conversation.ConversationID.Hex()))
		}
		seqs, err := l.svcCtx.Redis().MGet(l.ctx, seqKeys...).Result()
		if err != nil {
			l.Errorf("ListUserConversation failed, err: %v", err)
			return nil, err
		}
		for i, seq := range seqs {
			seqsMap[seqKeys[i]] = utils.InterfaceToInt64(seq)
		}
		unreadMap, err = l.svcCtx.Redis().HGetAll(l.ctx, unreadKey).Result()
		if err != nil {
			l.Errorf("ListUserConversation failed, err: %v", err)
			return nil, err
		}
		for _, conversation := range conversations {
			seq, ok := seqsMap[rediskey.ConversationSeq(conversation.ConversationID.Hex())]
			if !ok {
				seq = 0
			}
			unread, ok := unreadMap[conversation.ConversationID.Hex()]
			if !ok {
				unread = "0"
			}
			conversation.MaxSeq = uint32(seq)
			conversation.Unread = uint32(utils.InterfaceToInt64(unread))
		}
	}
	return conversations, nil
}

func (l *Default) PullMsgBySeq(conversationID string, userID string, maxSeq uint32, pageSize int64) ([]*model.ChatLog, error) {
	conversationId, err := primitive.ObjectIDFromHex(conversationID)
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(l.ctx, time.Second*time.Duration(l.svcCtx.Config.Mongo.DBTimeout))
	// 查询用户会话minseq
	var midUserConversation = &model.MidUserConversation{}
	filter := bsonx.Doc{}
	err = l.svcCtx.Collection(midUserConversation).FindOne(ctx,
		filter.Append("conversation_id", bsonx.ObjectID(conversationId)).Append("belong_user", bsonx.String(userID)).Append("deleted_at", bsonx.Int64(0)),
	).Decode(midUserConversation)
	if err != nil {
		l.Errorf("query mid user conversation error: %s", err.Error())
		return nil, err
	}
	cursor, err := l.svcCtx.Collection(&model.ChatLog{}).Find(ctx, bson.M{
		"conversation_id": conversationId,
		"seq": bson.M{
			"$gt": midUserConversation.MinSeq,
			"$lt": maxSeq,
		},
	}, options.Find().SetLimit(pageSize).SetSort(bson.M{"seq": 1}))
	if err != nil {
		l.Errorf("query chat log error: %s", err.Error())
		return nil, err
	}
	var chatLogs []*model.ChatLog
	err = cursor.All(ctx, &chatLogs)
	if err != nil {
		l.Errorf("query chat log error: %s", err.Error())
		return nil, err
	}
	return chatLogs, nil
}
