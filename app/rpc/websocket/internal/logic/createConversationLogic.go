package logic

import (
	"context"
	"errors"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/database"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/model"
	"github.com/IM-Lite/IM-Lite-Server/common/xtrace"
	"google.golang.org/protobuf/proto"
	"time"

	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/svc"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateConversationLogic {
	return &CreateConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateConversationLogic) CreateConversation(in *pb.CreateConversationReq) (*pb.CreateConversationResp, error) {
	if in.Members == nil {
		return &pb.CreateConversationResp{}, errors.New("members is nil")
	}
	ctx, _ := context.WithTimeout(l.ctx, time.Duration(l.svcCtx.Config.Mongo.DBTimeout)*time.Second)
	collection := &model.Conversation{
		Members: in.Members,
	}
	conversationID, err := database.NewDefault(l.svcCtx, ctx).CreateConversation(collection)
	if err != nil {
		l.Errorf("update conversation min seq error: %s", err.Error())
		return &pb.CreateConversationResp{}, err
	}
	{
		mq := &pb.IMMsgPushMQ{PushBody: &pb.PushBody{
			Event: pb.PushEvent_updateConv,
			Data:  nil,
		}, TraceId: xtrace.TraceIdFromContext(ctx)}
		value, _ := proto.Marshal(mq)
		ctx, _ = context.WithTimeout(ctx, time.Second*1)
		_, _, err = l.svcCtx.PushPusher().SendMessage(ctx, value, conversationID)
		if err != nil {
			logx.Errorf("push to kafka error, err: %s", err.Error())
		}
	}
	return &pb.CreateConversationResp{ConversationId: conversationID}, nil
}
