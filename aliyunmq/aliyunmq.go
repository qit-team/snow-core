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
