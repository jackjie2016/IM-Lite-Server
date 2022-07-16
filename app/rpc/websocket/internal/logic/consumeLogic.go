package logic

import (
	"context"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/database"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/model"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/svc"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/ws"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/pb"
	"github.com/IM-Lite/IM-Lite-Server/common/xkafka"
	"github.com/IM-Lite/IM-Lite-Server/common/xredis/rediskey"
	"github.com/IM-Lite/IM-Lite-Server/common/xtrace"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/proto"
	"time"
)

type ConsumeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func (l *ConsumeLogic) HandleMsg(value []byte, key string, topic string, partition int32, offset int64) error {
	switch topic {
	case l.svcCtx.Config.Consumers.Msg.Topic:
		msg := &pb.IMMsgDataMQ{}
		err := proto.Unmarshal(value, msg)
		if err != nil {
			l.Errorf("unmarshal msg error: %s", err.Error())
			return nil
		}
		if msg.Msg == nil {
			l.Errorf("msg is nil")
			return nil
		}
		xtrace.RunWithTrace(msg.TraceId, "websocket.consume.store", func(ctx context.Context) {
			err = l.ConsumerMsg(ctx, key, msg, value)
		})
		return err
	case l.svcCtx.Config.Consumers.Push.Topic:
		msg := &pb.IMMsgPushMQ{}
		err := proto.Unmarshal(value, msg)
		if err != nil {
			l.Errorf("unmarshal msg error: %s", err.Error())
			return nil
		}
		if msg.PushBody == nil {
			l.Errorf("PushBody is nil")
			return nil
		}
		xtrace.RunWithTrace(msg.TraceId, "websocket.consume.push", func(ctx context.Context) {
			l.ConsumerPush(ctx, key, msg)
		})
		return nil
	}
	return nil
}

func (l *ConsumeLogic) ConsumerMsg(ctx context.Context, conversationIDStr string, msg *pb.IMMsgDataMQ, value []byte) error {
	logger := l.WithContext(ctx)
	// 新增消息
	conversationID, err := primitive.ObjectIDFromHex(conversationIDStr)
	if err != nil {
		logger.Errorf("conversationIDStr to ObjectID error, conversationID: %s, err: %s", conversationIDStr, err)
		return nil
	}
	{ // 是否有这个会话
		has, err := database.NewDefault(l.svcCtx, l.ctx).HasConversation(conversationIDStr)
		if err != nil {
			l.Errorf("has conversation error: %s", err.Error())
			return err
		}
		if !has {
			logger.Errorf("has conversation error, conversationID: %s", conversationIDStr)
			return nil
		}
	}
	var (
		now     = time.Now().UnixNano()
		chatLog = &model.ChatLog{
			ConversationID: conversationID,
			ClientMsgID:    msg.Msg.ClientMsgID,
			ClientTime:     msg.Msg.ClientTime,
			ServerTime:     now,
			SenderID:       msg.Msg.SenderID,
			Seq:            0,
			Data: model.ChatLogData{
				ContentType: msg.Msg.ContentType,
				Content:     msg.Msg.Content,
			},
		}
	)
	if msg.Msg.IsStore() {
		err = database.NewDefault(l.svcCtx, ctx).InsertMayUpdateOneMessage(chatLog)
		if err != nil {
			l.Errorf("insert message error: %s", err.Error())
			return err
		}
	}
	// 推送
	msgData := &pb.MsgData{
		ClientMsgID: chatLog.ClientMsgID,
		ServerMsgID: chatLog.ServerMsgID.Hex(),
		SenderID:    chatLog.SenderID,
		ConvID:      chatLog.ConversationID.Hex(),
		ContentType: chatLog.Data.ContentType,
		Content:     chatLog.Data.Content,
		ClientTime:  chatLog.ClientTime,
		ServerTime:  chatLog.ServerTime,
		Seq:         chatLog.Seq,
		OfflinePush: msg.Msg.OfflinePush,
		MsgOptions:  msg.Msg.MsgOptions,
	}
	data, _ := proto.Marshal(msgData)
	mq := &pb.IMMsgPushMQ{PushBody: &pb.PushBody{
		Event: pb.PushEvent_receiveMsg,
		Data:  data,
	}, TraceId: xtrace.TraceIdFromContext(ctx)}
	value, _ = proto.Marshal(mq)
	ctx, _ = context.WithTimeout(ctx, time.Second*1)
	_, _, err = l.svcCtx.PushPusher().SendMessage(ctx, value, conversationIDStr)
	if err != nil {
		logx.Errorf("push to kafka error, err: %s", err.Error())
	}
	return nil
}

func (l *ConsumeLogic) ConsumerPush(ctx context.Context, conversationIDStr string, msg *pb.IMMsgPushMQ) {
	logger := l.WithContext(ctx)
	listAllSubscribersByConversationResp, err := NewListAllSubscribersByConversationLogic(ctx, l.svcCtx).ListAllSubscribersByConversation(&pb.ListAllSubscribersByConversationReq{ConversationId: conversationIDStr})
	if err != nil {
		logger.Errorf("list all subscribers by conversation error: %s", err.Error())
		return
	}
	var msgData *pb.MsgData
	body := msg.PushBody
	switch body.Event {
	case pb.PushEvent_receiveMsg:
		msgData = &pb.MsgData{}
		err = proto.Unmarshal(body.Data, msgData)
		if err != nil {
			logger.Errorf("unmarshal msg error: %s", err.Error())
			return
		}
	}
	w := ws.GetWs(l.svcCtx)
	for _, uid := range listAllSubscribersByConversationResp.UserIds {
		if msgData != nil && msgData.IsIncrUnread() {
			err := l.svcCtx.Redis().HIncrBy(ctx, rediskey.ConversationUnread(uid), conversationIDStr, 1).Err()
			if err != nil {
				logger.Errorf("incr conversation unread error: %s", err.Error())
			}
		}
		bodyValue, _ := proto.Marshal(body)
		w.SendMsgToUidIgnoreErr(ctx, uid, bodyValue)
	}
}

func (l *ConsumeLogic) run() {
	msg := xkafka.GetConsumer(l.svcCtx.Config.GetConsumersMsg())
	go msg.RegisterHandleAndConsumer(l)
	push := xkafka.GetConsumer(l.svcCtx.Config.GetConsumersPush())
	go push.RegisterHandleAndConsumer(l)
}

func NewConsumeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConsumeLogic {
	c := &ConsumeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
	c.run()
	return c
}
