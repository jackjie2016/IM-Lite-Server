package logic

import (
	"context"
	"github.com/IM-Lite/IM-Lite-Server/common/xtrace"
	"google.golang.org/protobuf/proto"
	"time"

	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/svc"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMsgLogic {
	return &SendMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendMsgLogic) SendMsg(in *pb.SendMsgReq) (*pb.SendMsgResp, error) {
	data := &pb.IMMsgDataMQ{TraceId: xtrace.TraceIdFromContext(l.ctx), Msg: in.Msg}
	v, _ := proto.Marshal(data)
	ctx, _ := context.WithTimeout(l.ctx, time.Second*1)
	_, _, err := l.svcCtx.MsgPusher().SendMessage(ctx, v, in.Msg.ConvID)
	if err != nil {
		l.Errorf("push msg failed, err: %v", err)
		return &pb.SendMsgResp{FailedMsg: "服务繁忙"}, err
	}
	return &pb.SendMsgResp{}, nil
}
