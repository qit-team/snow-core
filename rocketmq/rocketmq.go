package rocketmq

import (
	"errors"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/qit-team/snow-core/config"
)

//依赖注入用的函数
func NewRocketMqClient(mqConfig config.RocketMqConfig) (client *RocketClient, err error) {
	// 初始化aliyunmq的 client
	defer func() {
		if e := recover(); e != nil {
			s := fmt.Sprintf("rocketmq client init panic: %v", e)
			err = errors.New(s)
		}
	}()

	if mqConfig.EndPoint != "" {
		client = new(RocketClient)

		consumerOptions := make([]consumer.Option, 0)
		consumerOptions = []consumer.Option{
			consumer.WithNameServer([]string{mqConfig.EndPoint}),
			consumer.WithCredentials(primitive.Credentials{
				AccessKey: mqConfig.AccessKey,
				SecretKey: mqConfig.SecretKey,
			}),
			consumer.WithGroupName(mqConfig.GroupId),
			consumer.WithNamespace(mqConfig.InstanceId),
			consumer.WithConsumerModel(consumer.Clustering),
			consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset),
			consumer.WithAutoCommit(true),
		}
		if len(mqConfig.ConsumerOptions) > 0 {
			consumerOptions = append(consumerOptions, mqConfig.ConsumerOptions...)
		}

		client.Consumer, err = rocketmq.NewPushConsumer(consumerOptions...)
		if err != nil {
			return nil, err
		}

		producerOptions := make([]producer.Option, 0)
		producerOptions = []producer.Option{
			producer.WithNameServer([]string{mqConfig.EndPoint}),
			producer.WithCredentials(primitive.Credentials{
				AccessKey: mqConfig.AccessKey,
				SecretKey: mqConfig.SecretKey,
			}),
			producer.WithRetry(2),
			producer.WithGroupName(mqConfig.GroupId),
			producer.WithNamespace(mqConfig.InstanceId),
		}
		if len(mqConfig.ProducerOptions) > 0 {
			producerOptions = append(producerOptions, mqConfig.ProducerOptions...)
		}

		client.Producer, err = rocketmq.NewProducer(producerOptions...)
		if err != nil {
			return nil, err
		}

	} else {
		err = errors.New("EndPoint empty, can not get client")
	}
	return
}

type RocketClient struct {
	Consumer rocketmq.PushConsumer
	Producer rocketmq.Producer
}
