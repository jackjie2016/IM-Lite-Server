package handler

import (
	"github.com/IM-Lite/IM-Lite-Server/common/xhttp"
	"io/ioutil"
	"net/http"

	"github.com/IM-Lite/IM-Lite-Server/app/api/internal/logic"
	"github.com/IM-Lite/IM-Lite-Server/app/api/internal/svc"
	"github.com/IM-Lite/IM-Lite-Server/app/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func PullMessagesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		bytes, e := ioutil.ReadAll(r.Body)
		if e != nil {
			httpx.OkJson(w, xhttp.NewParamErr(e))
			return
		}
		var req types.ReqPullMessages
		req.Message = bytes
		l := logic.NewPullMessagesLogic(r.Context(), svcCtx)
		resp, err := l.PullMessages(&req)
		if err != nil {
			httpx.OkJson(w, xhttp.Failed(err))
		} else {
			_, _ = w.Write(resp.MessageLists)
			httpx.Ok(w)
		}
	}
}
