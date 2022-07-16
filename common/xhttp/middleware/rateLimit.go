package middleware

import (
	"context"
	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/logx"
)

func (m *AuthMiddleware) RateLimit(ctx context.Context, ip string) bool {
	logger := logx.WithContext(ctx)
	takeRes, err := m.ipPeriodLimit.TakeCtx(ctx, ip)
	if err != nil {
		logger.Errorf("ip:%s, rate limit err:", err)
		return true
	}
	switch takeRes {
	case limit.OverQuota:
		logger.Errorf("ip:%s, rate limit OverQuota", ip)
		return false
	case limit.Allowed:
		return true
	case limit.HitQuota:
		logger.Errorf("ip:%s, rate limit HitQuota", ip)
		return false
	default:
		return true
	}
}
