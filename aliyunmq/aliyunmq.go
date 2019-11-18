package aliyunmq

import (
	"errors"
	"fmt"
	"github.com/aliyunmq/mq-http-go-sdk"
	"github.com/qit-team/snow-core/config"
)

//依赖注入用的函数
func NewAliyunMqClient(mqConfig config.AliyunMqConfig) (client mq_http_sdk.MQClient, err error) {
	// 初始化aliyunmq的 client
	defer func() {
		if e := recover(); e != nil {
			s := fmt.Sprintf("aliyun_mq client init panic: %s", fmt.Sprint(e))
			err = errors.New(s)
		}
	}()

	if mqConfig.EndPoint != "" {
		client = mq_http_sdk.NewAliyunMQClient(mqConfig.EndPoint, mqConfig.AccessKey, mqConfig.SecretKey, "")
	} else {
		err = errors.New("EndPoint empty,can not get client")
	}
	return
}

// func GetAliyunMqProducer(client mq_http_sdk.MQClient, )
//
//func GetMnsBasicQueue(client ali_mns.MNSClient, queueName string) ali_mns.AliMNSQueue {
//	var defaultQueue ali_mns.AliMNSQueue
//
//	//根据client创建manager
//	queueManager := ali_mns.NewMNSQueueManager(client)
//
//	// 暂时将visibilityTimeout 设置成120，后续将参数暴露给上层，可自行配置
//	err := queueManager.CreateQueue(queueName, 0, 65536, 345600, 120, 0, 3)
//	if err != nil && !ali_mns.ERR_MNS_QUEUE_ALREADY_EXIST_AND_HAVE_SAME_ATTR.IsEqual(err) {
//		fmt.Println(err)
//		return defaultQueue
//	}
//	//最终的最小执行单元queue
//	return ali_mns.NewMNSQueue(queueName, client)
//}
