package logic

import (
	"context"
	"github.com/IM-Lite/IM-Lite-Server/common/xredis/rediskey"

	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/svc"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClearUnreadCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClearUnreadCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearUnreadCountLogic {
	return &ClearUnreadCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ClearUnreadCountLogic) ClearUnreadCount(in *pb.ClearUnreadCountReq) (*pb.ClearUnreadCountResp, error) {
	var err error
	if in.ConversationId == "" {
		err = l.svcCtx.Redis().Del(l.ctx, rediskey.ConversationUnread(in.UserId)).Err()
	} else {
		err = l.svcCtx.Redis().HDel(l.ctx, rediskey.ConversationUnread(in.UserId), in.ConversationId).Err()
	}
	if err != nil {
		l.Errorf("clear unread count err: %v", err)
		return &pb.ClearUnreadCountResp{}, err
	}
	return &pb.ClearUnreadCountResp{}, nil
}
