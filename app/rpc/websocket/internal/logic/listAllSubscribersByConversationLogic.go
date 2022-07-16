package logic

import (
	"context"
	"github.com/IM-Lite/IM-Lite-Server/common/utils"
	"github.com/IM-Lite/IM-Lite-Server/common/xredis/rediskey"
	"strings"
	"time"

	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/svc"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAllSubscribersByConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAllSubscribersByConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAllSubscribersByConversationLogic {
	return &ListAllSubscribersByConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListAllSubscribersByConversationLogic) ListAllSubscribersByConversation(in *pb.ListAllSubscribersByConversationReq) (*pb.ListAllSubscribersByConversationResp, error) {
	now := time.Now().UnixMilli()
	result, _, err := l.svcCtx.Redis().Scan(l.ctx, 0, rediskey.UserSubscribedConversations(in.ConversationId, "*"), 10000).Result()
	if err != nil {
		l.Errorf("list all subscribers by conversation failed, err: %v", err)
		return nil, err
	}
	if len(result) == 0 {
		return &pb.ListAllSubscribersByConversationResp{}, nil
	}
	values, err := l.svcCtx.Redis().MGet(l.ctx, result...).Result()
	if err != nil {
		l.Errorf("list all subscribers by conversation failed, err: %v", err)
	} else {
		var uids []string
		for i, v := range result {
			value := values[i]
			timestamp := utils.InterfaceToInt64(value)
			if now-timestamp < 60*1000 {
				uids = append(uids, strings.TrimPrefix(v, rediskey.UserSubscribedConversations(in.ConversationId, "")))
			}
		}
		return &pb.ListAllSubscribersByConversationResp{UserIds: uids}, nil
	}
	return &pb.ListAllSubscribersByConversationResp{}, nil
}
