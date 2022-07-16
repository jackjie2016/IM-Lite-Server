package middleware

import (
	"context"
	"github.com/IM-Lite/IM-Lite-Server/common/utils"
	"github.com/IM-Lite/IM-Lite-Server/common/xhttp"
	"github.com/IM-Lite/IM-Lite-Server/common/xredis"
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"strings"
)

type IGetCache interface {
	Cache() redis.UniversalClient
}
type AuthMiddleware struct {
	universalClient redis.UniversalClient
	ipPeriodLimit   *limit.PeriodLimit
}

func NewAuthMiddleware(c IGetCache, config xredis.Config) *AuthMiddleware {
	universalClient := c.Cache()
	redisClient := xredis.GetZeroRedis(config)
	ipPeriodLimit := limit.NewPeriodLimit(10, 10, redisClient, "middleware:periodlimit:notoken:ip:")
	return &AuthMiddleware{universalClient: universalClient, ipPeriodLimit: ipPeriodLimit}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("token")
		ip := req.Header.Get("x-real-ip")
		ua := req.Header.Get("user-agent")
		av := req.Header.Get("app-version")
		platform := req.Header.Get("platform")
		userId := req.Header.Get("User_id")
		logger := logx.WithContext(req.Context())
		if len(strings.Split(token, ".")) != 3 {
			// 判断是否是白名单接口
			if m.isWhite(req) {
				ctx := context.WithValue(req.Context(), xhttp.CtxKeyUserId, userId)
				ctx = context.WithValue(ctx, xhttp.CtxKeyPlatform, platform)
				// 根据ip限流
				if !m.RateLimit(ctx, ip) {
					logger.Error("[业务告警]", "用户频繁请求白名单接口：", "ip：", ip, "ua：", ua, "av：", av, "uri：", req.RequestURI)
					httpx.OkJson(resp, xhttp.Failed(xhttp.NewErr(xhttp.IpRateLimitErrCode, "您操作太过频繁")))
					return
				}
				next(resp, req.WithContext(ctx))
				return
			}
			httpx.OkJson(resp, xhttp.Failed(xhttp.NewErr(xhttp.AuthErrCode, "登录已过期")))
			logger.Error("[业务告警]", "用户破解接口访问失败：", "token：", token, "ip：", ip, "ua：", ua, "av：", av)
			return
		}
		uid, err := m.VerifyToken(req.Context(), token, platform, userId)
		if err != nil {
			httpx.OkJson(resp, xhttp.Failed(err))
			logger.Error("[业务告警]", "用户破解接口访问失败：", "token：", token, "ip：", ip, "ua：", ua, "av：", av)
			return
		}
		{
			req.Header.Set("user-id", uid)
			ctx := context.WithValue(req.Context(), xhttp.CtxKeyUserId, uid)
			ctx = context.WithValue(ctx, xhttp.CtxKeyPlatform, platform)
			ctx = context.WithValue(ctx, xhttp.CtxKeyAppVersion, av)
			ctx = context.WithValue(ctx, xhttp.CtxKeyUserAgent, ua)
			ctx = context.WithValue(ctx, xhttp.CtxKeyIp, ip)
			ctx = context.WithValue(ctx, xhttp.CtxKeyToken, token)
			next(resp, req.WithContext(ctx))
		}
		return
	}
}

func (m *AuthMiddleware) VerifyToken(ctx context.Context, token string, platform string, id string) (string, xhttp.ICodeErr) {
	status, uid := utils.VerifyToken(ctx, m.universalClient, token)
	switch status {
	case utils.VerifyTokenCodeBaned:
		return "", xhttp.NewErr(xhttp.AuthBanedErrCode, "您的账号已被禁用")
	case utils.VerifyTokenCodeExpire:
		return "", xhttp.NewErr(xhttp.AuthErrCode, "登录已过期")
	case utils.VerifyTokenCodeOK:
		return uid, nil
	}
	return "", xhttp.NewDefaultErr("服务繁忙，请稍后再试")
}
