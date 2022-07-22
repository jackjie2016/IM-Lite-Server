package logic

import (
	"context"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/database"

	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/svc"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PullMsgBySeqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPullMsgBySeqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PullMsgBySeqLogic {
	return &PullMsgBySeqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PullMsgBySeqLogic) PullMsgBySeq(in *pb.PullMsgBySeqReq) (*pb.PullMsgBySeqResp, error) {
	chatLogs, err := database.NewDefault(l.svcCtx, l.ctx).PullMsgBySeq(in.ConversationId, in.UserId, in.OldestSeq, in.GetPageSizeX())
	if err != nil {
		l.Errorf("PullMsgBySeq failed, err: %v", err)
		return &pb.PullMsgBySeqResp{}, err
	}
	var resp []*pb.MsgData
	for _, chatLog := range chatLogs {
		resp = append(resp, &pb.MsgData{
			ClientMsgID: chatLog.ClientMsgID,
			ServerMsgID: chatLog.ServerMsgID.Hex(),
			SenderID:    chatLog.SenderID,
			ConvID:      chatLog.ConversationID.Hex(),
			ContentType: chatLog.Data.ContentType,
			Content:     chatLog.Data.Content,
			ClientTime:  chatLog.ClientTime,
			ServerTime:  chatLog.ServerTime,
			Seq:         chatLog.Seq,
			OfflinePush: &pb.OfflinePush{
				Title:         chatLog.OfflinePush.Title,
				Desc:          chatLog.OfflinePush.Desc,
				Ex:            chatLog.OfflinePush.Ex,
				IOSPushSound:  chatLog.OfflinePush.IOSPushSound,
				IOSBadgeCount: chatLog.OfflinePush.IOSBadgeCount,
			},
			MsgOptions: &pb.MsgOptions{
				Storage: chatLog.MsgOptions.Storage,
				Unread:  chatLog.MsgOptions.Unread,
			},
		})
	}
	return &pb.PullMsgBySeqResp{MsgList: resp, ConversationId: in.ConversationId}, nil
}
