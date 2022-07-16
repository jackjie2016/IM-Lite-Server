package logic

import (
	"context"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/database"

	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/svc"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteConversationAllMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteConversationAllMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteConversationAllMsgLogic {
	return &DeleteConversationAllMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteConversationAllMsgLogic) DeleteConversationAllMsg(in *pb.DeleteConversationAllMsgReq) (*pb.DeleteConversationAllMsgResp, error) {
	err := database.NewDefault(l.svcCtx, l.ctx).DeleteConversationAllMessages(in.ConversationId, in.UserId)
	if err != nil {
		l.Errorf("DeleteConversationAllMsg error: %v", err)
	}
	return &pb.DeleteConversationAllMsgResp{}, err
}
