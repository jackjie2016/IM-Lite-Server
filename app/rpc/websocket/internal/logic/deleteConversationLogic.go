package logic

import (
	"context"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/database"
	"github.com/IM-Lite/IM-Lite-Server/common/xtrace"
	"google.golang.org/protobuf/proto"
	"time"

	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/svc"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteConversationLogic {
	return &DeleteConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteConversationLogic) DeleteConversation(in *pb.DeleteConversationReq) (*pb.DeleteConversationResp, error) {
	err := database.NewDefault(l.svcCtx, l.ctx).DeleteConversation(in.ConversationId)
	if err != nil {
		l.Errorf("DeleteConversation error: %v", err)
	}
	{
		mq := &pb.IMMsgPushMQ{PushBody: &pb.PushBody{
			Event: pb.PushEvent_updateConv,
			Data:  nil,
		}, TraceId: xtrace.TraceIdFromContext(l.ctx)}
		value, _ := proto.Marshal(mq)
		ctx, _ := context.WithTimeout(l.ctx, time.Second*1)
		_, _, err = l.svcCtx.PushPusher().SendMessage(ctx, value, in.ConversationId)
		if err != nil {
			logx.Errorf("push to kafka error, err: %s", err.Error())
		}
	}
	return &pb.DeleteConversationResp{}, err
}
