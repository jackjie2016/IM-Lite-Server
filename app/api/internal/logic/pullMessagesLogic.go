package logic

import (
	"context"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/pb"
	"github.com/IM-Lite/IM-Lite-Server/common/xhttp"
	"github.com/zeromicro/go-zero/core/mr"
	"google.golang.org/protobuf/proto"

	"github.com/IM-Lite/IM-Lite-Server/app/api/internal/svc"
	"github.com/IM-Lite/IM-Lite-Server/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PullMessagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPullMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PullMessagesLogic {
	return &PullMessagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PullMessagesLogic) PullMessages(req *types.ReqPullMessages) (resp *types.RespPullMessages, ierr xhttp.ICodeErr) {
	in := &pb.PullMsgList{}
	err := proto.Unmarshal(req.Message, in)
	if err != nil {
		return nil, xhttp.NewParamErr(err)
	}

	var messageList = &pb.MsgDataList{}
	uid := xhttp.GetUidFromCtx(l.ctx)
	type Conv struct {
		ConvId    string
		SeqList   []uint32
		OldestSeq uint32
		PageSize  int32
	}
	var convs []*Conv
	for _, conv := range in.List {
		// 验证参数
		if len(conv.ConvID) == 0 {
			return nil, xhttp.NewParamErrByMsg("conv_id不能为空")
		}
		// conv.SeqList 是否是连续的 从小到大的
		for i, u := range conv.SeqList {
			if i == 0 {
				continue
			}
			if u+1 != conv.SeqList[i-1] {
				return nil, xhttp.NewParamErrByMsg("seqList不是连续加一的")
			}
		}
		convs = append(convs, &Conv{
			ConvId:    conv.ConvID,
			SeqList:   conv.SeqList,
			OldestSeq: conv.SeqList[len(conv.SeqList)-1] + 1,
			PageSize:  int32(len(conv.SeqList)),
		})
	}
	var fs []func() error
	for _, conv := range convs {
		fs = append(fs, func() error {
			pullMsgBySeqResp, err := l.svcCtx.WebsocketService().PullMsgBySeq(l.ctx, &pb.PullMsgBySeqReq{
				ConversationId: conv.ConvId,
				UserId:         uid,
				OldestSeq:      conv.OldestSeq,
				PageSize:       conv.PageSize,
			})
			if err != nil {
				l.Errorf("PullMsgBySeq err: %v", err)
				return nil
			}
			msgList := pullMsgBySeqResp.MsgList
			var seqMap = make(map[uint32]*pb.MsgData)
			for _, data := range msgList {
				if _, ok := seqMap[data.Seq]; ok {
					l.Errorf("PullMsgBySeq err: seq重复")
					return nil
				}
				seqMap[data.Seq] = data
			}
			for _, reqSeq := range conv.SeqList {
				if data, ok := seqMap[reqSeq]; !ok {
					l.Info("PullMsgBySeq err: seq不存在")
					messageList.List = append(messageList.List, &pb.MsgData{Seq: reqSeq})
				} else {
					messageList.List = append(messageList.List, data)
				}
			}
			return nil
		})
	}
	_ = mr.Finish(fs...)
	buf, _ := proto.Marshal(messageList)
	resp = &types.RespPullMessages{
		MessageLists: buf,
	}
	return
}
