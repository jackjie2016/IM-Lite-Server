package xkafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

type ConsumerGroup struct {
	*kafka.Consumer
	groupID string
	topics  []string
	cfg     ConsumerGroupConfig
}
type msgConsumer interface {
	HandleMsg(value []byte, key string, topic string, partition int32, offset int64) error
}

type ConsumerGroupConfig struct {
	Addrs     []string
	Topic     string
	Group     string
	Offset    string `json:",options=first|last,default=last"`
	Consumers int    `json:",default=8"`
	User      string `json:",optional"`
	Passwd    string `json:",optional"`
}

func (c ConsumerGroupConfig) GetOffset() kafka.ConfigValue {
	if c.Offset == "first" {
		return "earliest"
	}
	return "latest"
}

func GetConsumer(
	cfg ConsumerGroupConfig,
) *ConsumerGroup {
	conf := &kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(cfg.Addrs, ","),
		"auto.offset.reset":        cfg.GetOffset(),
		"enable.auto.offset.store": false,
		"enable.auto.commit":       true,
		"retries":                  10000000,
		"group.id":                 cfg.Group,
	}
	if cfg.User != "" {
		_ = conf.Set("sasl.mechanisms=PLAIN")
		_ = conf.Set("security.protocol=SASL_PLAINTEXT")
		_ = conf.Set("sasl.username=" + cfg.User)
		_ = conf.Set("sasl.password=" + cfg.Passwd)
	}
	c, err := kafka.NewConsumer(conf)
	if err != nil {
		logx.Errorf("new consumer error: %s", err.Error())
		panic(err.Error())
	}
	return &ConsumerGroup{
		c,
		cfg.Group,
		[]string{cfg.Topic},
		cfg,
	}
}

func (mc *ConsumerGroup) RegisterHandleAndConsumer(handler msgConsumer) {
	err := mc.Consumer.SubscribeTopics(mc.topics, nil)
	if err != nil {
		panic(err.Error())
	}
	for i := 1; i < mc.cfg.Consumers; i++ {
		go mc.loop(handler)
	}
	mc.loop(handler)
}

func (mc *ConsumerGroup) loop(handler msgConsumer) {
	for {
		ev := mc.Consumer.Poll(100)
		if ev == nil {
			continue
		}
		switch e := ev.(type) {
		case *kafka.Message:
			err1 := handler.HandleMsg(e.Value, string(e.Key), *e.TopicPartition.Topic, e.TopicPartition.Partition, int64(e.TopicPartition.Offset))
			if err1 == nil {
				if _, er := mc.Consumer.StoreMessage(e); er != nil {
					logx.Errorf("store message error: %s", er.Error())
				}
			}
		}
	}
}
