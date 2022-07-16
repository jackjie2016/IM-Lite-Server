package logic

import (
	"context"
	"github.com/IM-Lite/IM-Lite-Server/common/xredis/rediskey"
	"time"

	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/svc"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSubscribedConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateSubscribedConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSubscribedConversationsLogic {
	return &UpdateSubscribedConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateSubscribedConversationsLogic) UpdateSubscribedConversations(in *pb.UpdateSubscribedConversationsReq) (*pb.UpdateSubscribedConversationsResp, error) {
	// 每100个切一片
	for i := 0; i < len(in.ConversationIds); i += 100 {
		end := i + 100
		if end > len(in.ConversationIds) {
			end = len(in.ConversationIds)
		}
		var kvs []interface{}
		now := time.Now().UnixMilli()
		for _, v := range in.ConversationIds[i:end] {
			// 用户最后活跃时间 用户最后订阅会话时间 now
			kvs = append(kvs, rediskey.UserSubscribedConversations(v, in.UserId), now)
		}
		err := l.svcCtx.Redis().MSet(l.ctx, kvs...).Err()
		if err != nil {
			l.Errorf("update subscribed conversations failed, err: %v", err)
			return nil, err
		}
	}
	return &pb.UpdateSubscribedConversationsResp{}, nil
}
