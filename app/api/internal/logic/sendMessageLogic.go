package logic

import (
	"context"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/pb"
	"github.com/IM-Lite/IM-Lite-Server/common/xhttp"
	"google.golang.org/protobuf/proto"

	"github.com/IM-Lite/IM-Lite-Server/app/api/internal/svc"
	"github.com/IM-Lite/IM-Lite-Server/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendMessageLogic) SendMessage(req *types.ReqSendMessage) (resp *types.RespSendMessage, ierr xhttp.ICodeErr) {
	msgData := &pb.MsgData{}
	e := proto.Unmarshal(req.Message, msgData)
	if e != nil {
		l.Errorf("unmarshal msgData error: %v", e)
		return nil, xhttp.NewParamErr(e)
	}
	sendMsgResp, err := l.svcCtx.WebsocketService().SendMsg(l.ctx, &pb.SendMsgReq{
		Msg: msgData,
	})
	if err != nil {
		l.Errorf("SendMsg err: %v", err)
		return nil, xhttp.NewInternalErr()
	}
	resp = &types.RespSendMessage{
		FailedMsg: sendMsgResp.FailedMsg,
	}
	return
}
