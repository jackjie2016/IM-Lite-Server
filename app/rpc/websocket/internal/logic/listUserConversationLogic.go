package logic

import (
	"context"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/database"

	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/svc"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUserConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListUserConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserConversationLogic {
	return &ListUserConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListUserConversationLogic) ListUserConversation(in *pb.ListUserConversationReq) (*pb.ListUserConversationResp, error) {
	var (
		resp = &pb.ListUserConversationResp{}
	)
	userConversations, err := database.NewDefault(l.svcCtx, l.ctx).ListUserConversations(in.UserId, int(in.PageNo), int64(in.PageSize))
	if err != nil {
		l.Errorf("ListUserConversation failed, err: %v", err)
		return &pb.ListUserConversationResp{}, err
	}
	// 订阅用户的会话
	{
		go func() {
			var convIds []string
			for _, conversation := range userConversations {
				convIds = append(convIds, conversation.ConversationID.Hex())
			}
			_, e := NewUpdateSubscribedConversationsLogic(context.Background(), l.svcCtx).UpdateSubscribedConversations(&pb.UpdateSubscribedConversationsReq{
				UserId:          in.UserId,
				ConversationIds: convIds,
			})
			if e != nil {
				l.Errorf("ListUserConversation failed, err: %v", e)
			}
		}()
	}
	{
		for _, conversation := range userConversations {
			resp.UserConversations = append(resp.UserConversations, &pb.UserConversation{
				Id:         conversation.ConversationID.Hex(),
				Name:       conversation.Remark,
				Avatar:     conversation.Avatar,
				Seq:        conversation.MaxSeq,
				MinSeq:     conversation.MinSeq,
				Unread:     conversation.Unread,
				Timestamp:  int64(conversation.Timestamp.T),
				Top:        conversation.Top,
				MsgRecvOpt: conversation.RecvOpt,
				IsDeleted:  false,
			})
		}
	}
	return resp, nil
}
