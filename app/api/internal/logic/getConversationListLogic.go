package logic

import (
	"context"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/pb"
	"github.com/IM-Lite/IM-Lite-Server/common/xhttp"

	"github.com/IM-Lite/IM-Lite-Server/app/api/internal/svc"
	"github.com/IM-Lite/IM-Lite-Server/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetConversationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationListLogic {
	return &GetConversationListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetConversationListLogic) GetConversationList(req *types.ReqGetConversationList) (resp *types.RespGetConversationList, ierr xhttp.ICodeErr) {
	listUserConversationResp, e := l.svcCtx.WebsocketService().ListUserConversation(l.ctx, &pb.ListUserConversationReq{
		UserId: xhttp.GetUidFromCtx(l.ctx),
	})
	if e != nil {
		l.Errorf("ListUserConversation error: %v", e)
		return &types.RespGetConversationList{}, xhttp.NewInternalErr()
	} else {
		l.Infof("ListUserConversation resp: %v", listUserConversationResp.String())
	}
	resp = &types.RespGetConversationList{}
	for _, conversation := range listUserConversationResp.UserConversations {
		resp.Conversations = append(resp.Conversations, &types.ModelConversation{
			ConvId:      conversation.Id,
			MaxSeq:      conversation.Seq,
			MinSeq:      conversation.MinSeq,
			UnreadCount: conversation.Unread,
		})
	}
	return
}
