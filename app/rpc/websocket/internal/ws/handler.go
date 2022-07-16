package ws

import (
	"context"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/svc"
	"github.com/IM-Lite/IM-Lite-Server/common/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func Handler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ConnRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		websocket := GetWs(svcCtx)
		resp, ok := requestClient(r.Context(), svcCtx, &req)
		status := http.StatusUnauthorized
		if ok {
			err := websocket.WsUpgrade(r.Context(), resp.Uid, &req, w, r)
			if err != nil {
				logx.WithContext(r.Context()).Errorf("WsUpgrade error: %s", err)
				return
			}
		} else {
			w.Header().Set("Sec-Websocket-Version", "13")
			w.Header().Set("ws_err_msg", "args err, need token, sendID, platformID")
			http.Error(w, http.StatusText(status), status)
		}
	}
}

func requestClient(ctx context.Context, svcCtx *svc.ServiceContext, req *ConnRequest) (*ConnResponse, bool) {
	var resUid = req.UserID
	if req.UserID == "" || req.Platform == "" {
		return &ConnResponse{
			Uid:     "",
			ErrMsg:  "参数错误",
			Success: false,
		}, false
	}
	code, _ := utils.VerifyToken(ctx, svcCtx.Redis(), req.Token)
	if code != utils.VerifyTokenCodeOK {
		return &ConnResponse{
			Uid:     "",
			ErrMsg:  "token错误",
			Success: false,
		}, false
	}
	//resUid = uid
	return &ConnResponse{
		Uid:     resUid,
		ErrMsg:  "",
		Success: true,
	}, true
}
