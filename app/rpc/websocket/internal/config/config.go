package config

import (
	"github.com/IM-Lite/IM-Lite-Server/common/xkafka"
	"github.com/IM-Lite/IM-Lite-Server/common/xmgo"
	"github.com/IM-Lite/IM-Lite-Server/common/xredis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Websocket struct {
		Port       int
		MaxConnNum int `json:",default=1000"`
	}
	Pushers struct {
		Msg  xkafka.ProducerConfig
		Push xkafka.ProducerConfig
	}
	Consumers struct {
		Msg  ConsumerConfig
		Push ConsumerConfig
	}

	Redis xredis.Config
	Mongo struct {
		xmgo.Config
		DBTimeout int    `json:",default=10"`
		Database  string `json:",default=im_lite"`
	}
}
type ConsumerConfig struct {
	Addrs     []string
	Topic     string
	Consumers int    `json:",default=8"`
	User      string `json:",optional"`
	Passwd    string `json:",optional"`
}

func (c *Config) PodIp() string {
	return PodIp
}

func (c *Config) GetConsumersPush() xkafka.ConsumerGroupConfig {
	group := c.ServiceConf.Name
	if PodName != "" {
		group = PodName
	}
	return xkafka.ConsumerGroupConfig{
		Addrs:     c.Consumers.Push.Addrs,
		Topic:     c.Consumers.Push.Topic,
		Group:     group,
		Offset:    "last",
		Consumers: c.Consumers.Push.Consumers,
		User:      c.Consumers.Push.User,
		Passwd:    c.Consumers.Push.Passwd,
	}
}

func (c *Config) GetConsumersMsg() xkafka.ConsumerGroupConfig {
	return xkafka.ConsumerGroupConfig{
		Addrs:     c.Consumers.Msg.Addrs,
		Topic:     c.Consumers.Msg.Topic,
		Group:     c.ServiceConf.Name,
		Offset:    "last",
		Consumers: c.Consumers.Msg.Consumers,
		User:      c.Consumers.Msg.User,
		Passwd:    c.Consumers.Msg.Passwd,
	}
}
