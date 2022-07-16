package logic

import (
	"context"
	"github.com/IM-Lite/IM-Lite-Server/common/xhttp"

	"github.com/IM-Lite/IM-Lite-Server/app/api/internal/svc"
	"github.com/IM-Lite/IM-Lite-Server/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClearHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClearHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearHistoryLogic {
	return &ClearHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClearHistoryLogic) ClearHistory(req *types.ReqClearHistory) (resp *types.RespClearHistory, ierr xhttp.ICodeErr) {
	// todo: add your logic here and delete this line

	return
}
