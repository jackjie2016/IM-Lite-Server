package svc

import (
	"github.com/IM-Lite/IM-Lite-Server/app/rpc/websocket/internal/config"
	"github.com/IM-Lite/IM-Lite-Server/common/xkafka"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

type ServiceContext struct {
	Config        config.Config
	redis         redis.UniversalClient
	mongo         *mongo.Client
	collectionMap sync.Map
	msgPusher     *xkafka.Producer
	pushPusher    *xkafka.Producer
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config: c,
	}
	s.Redis()
	s.Mongo()
	s.MsgPusher()
	s.PushPusher()
	return s
}
