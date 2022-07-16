package svc

import (
	"github.com/IM-Lite/IM-Lite-Server/common/xkafka"
	"github.com/IM-Lite/IM-Lite-Server/common/xmgo"
	"github.com/IM-Lite/IM-Lite-Server/common/xredis"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func (s *ServiceContext) MsgPusher() *xkafka.Producer {
	if s.msgPusher == nil {
		s.msgPusher = xkafka.MustNewProducer(s.Config.Pushers.Msg)
	}
	return s.msgPusher
}
func (s *ServiceContext) PushPusher() *xkafka.Producer {
	if s.pushPusher == nil {
		s.pushPusher = xkafka.MustNewProducer(s.Config.Pushers.Push)
	}
	return s.pushPusher
}

func (s *ServiceContext) Redis() redis.UniversalClient {
	if s.redis == nil {
		s.redis = xredis.GetClient(s.Config.Redis)
	}
	return s.redis
}

func (s *ServiceContext) Mongo() *mongo.Client {
	if s.mongo == nil {
		s.mongo = xmgo.GetClient(s.Config.Mongo.Config)
	}
	return s.mongo
}

func (s *ServiceContext) Collection(model xmgo.ICollection) *mongo.Collection {
	if c, ok := s.collectionMap.Load(model.CollectionName()); ok {
		return c.(*mongo.Collection)
	} else {
		c := s.Mongo().Database(s.Config.Mongo.Database).Collection(model.CollectionName())
		err := model.Indexes(c)
		if err != nil {
			log.Panicln(err)
		}
		s.collectionMap.Store(model.CollectionName(), c)
		return c
	}
}
