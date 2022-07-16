package handler

import (
	"net/http"

	"github.com/IM-Lite/IM-Lite-Server/app/api/internal/logic"
	"github.com/IM-Lite/IM-Lite-Server/app/api/internal/svc"
	"github.com/IM-Lite/IM-Lite-Server/app/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ClearUnreadCountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReqClearUnreadCount
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewClearUnreadCountLogic(r.Context(), svcCtx)
		resp, err := l.ClearUnreadCount(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
