package logic

import (
	"context"
	"github.com/IM-Lite/IM-Lite-Server/common/xhttp"

	"github.com/IM-Lite/IM-Lite-Server/app/api/internal/svc"
	"github.com/IM-Lite/IM-Lite-Server/app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClearUnreadCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClearUnreadCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearUnreadCountLogic {
	return &ClearUnreadCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClearUnreadCountLogic) ClearUnreadCount(req *types.ReqClearUnreadCount) (resp *types.RespClearUnreadCount, ierr xhttp.ICodeErr) {
	// todo: add your logic here and delete this line

	return
}
