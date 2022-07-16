package xkafka

import (
	"context"
	"github.com/IM-Lite/IM-Lite-Server/common/xtrace"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
	"strings"
)

type ProducerConfig struct {
	Addrs  []string `json:""`
	Topic  string   `json:""`
	User   string   `json:",optional"`
	Passwd string   `json:",optional"`
}

type Producer struct {
	topic         string
	addr          []string
	config        ProducerConfig
	producer      *kafka.Producer
	msgChan       chan kafka.Event
	partitionIDs  []int32
	nextPartition int
}

func MustNewProducer(config ProducerConfig) *Producer {
	conf := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(config.Addrs, ","),
	}
	if config.User != "" {
		_ = conf.Set("sasl.mechanisms=PLAIN")
		_ = conf.Set("sasl.username=" + config.User)
		_ = conf.Set("sasl.password=" + config.Passwd)
		// 超时时间
		_ = conf.Set("request.timeout.ms=500")
		// 重试
		_ = conf.Set("retries=0")
	}
	producer, err := kafka.NewProducer(conf)
	if err != nil {
		logx.Errorf("kafka.NewProducer error: %v", err)
		panic(err)
	}
	admin, err := kafka.NewAdminClientFromProducer(producer)
	if err != nil {
		logx.Errorf("kafka.NewAdminClientFromProducer error: %v", err)
		panic(err)
	}
	metadata, err := admin.GetMetadata(&config.Topic, false, 2000)
	if err != nil {
		logx.Errorf("kafka.GetMetadata error: %v", err)
		panic(err)
	}
	var partitionIDs []int32
	for i, partition := range metadata.Topics[config.Topic].Partitions {
		logx.Infof("kafka.GetMetadata: %s:%d partition: %+v", config.Topic, i, partition)
		partitionIDs = append(partitionIDs, partition.ID)
	}
	ch := make(chan kafka.Event)
	return &Producer{
		topic:         config.Topic,
		addr:          config.Addrs,
		config:        config,
		producer:      producer,
		msgChan:       ch,
		partitionIDs:  partitionIDs,
		nextPartition: 0,
	}
}

func (p *Producer) SendMessage(ctx context.Context, bMsg []byte, key string) (partition int32, offset int64, err error) {
	defer func() {
		p.nextPartition++
		if p.nextPartition >= len(p.partitionIDs) {
			p.nextPartition = 0
		}
	}()
	for i := 0; i < 3; i++ {
		select {
		case <-ctx.Done():
			return 0, 0, ctx.Err()
		default:
			xtrace.StartFuncSpan(ctx, "kafka.Producer.SendMessage:key:"+key+"retrycount:"+strconv.Itoa(i), func(ctx context.Context) {
				err = p.producer.Produce(&kafka.Message{TopicPartition: kafka.TopicPartition{
					Topic:     &p.topic,
					Partition: p.partitionIDs[p.nextPartition],
				}, Key: []byte(key), Value: bMsg}, p.msgChan)
			})
			if err != nil {
				continue
			}
			e := <-p.msgChan
			ev := e.(*kafka.Message)
			if ev.TopicPartition.Error != nil {
				err = ev.TopicPartition.Error
				continue
			}
			partition, offset = ev.TopicPartition.Partition, int64(ev.TopicPartition.Offset)
			return
		}
	}
	return
}
