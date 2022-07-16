package svc

import (
	"github.com/IM-Lite/IM-Lite-Server/app/api/internal/config"
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/websocketservice"
	"github.com/IM-Lite/IM-Lite-Server/common/xhttp/middleware"
	"github.com/IM-Lite/IM-Lite-Server/common/xredis"
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config           config.Config
	Auth             rest.Middleware
	redis            redis.UniversalClient
	websocketService websocketservice.WebsocketService
}

func (s *ServiceContext) Cache() redis.UniversalClient {
	if s.redis == nil {
		s.redis = xredis.GetClient(s.Config.Redis)
	}
	return s.redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config: c,
	}
	s.Auth = middleware.NewAuthMiddleware(s, s.Config.Redis).Handle
	return s
}
