package alirocketqueue

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	mq_http_sdk "github.com/aliyunmq/mq-http-go-sdk"
	"github.com/gogap/errors"
	"github.com/qit-team/snow-core/aliyunmq"
	"github.com/qit-team/snow-core/queue"
)

const (
	DefaultVisibilityTimeout = int64(120)
)

var (
	mp map[string]queue.Queue
	mu sync.RWMutex
)

type AliyunMq struct {
	client mq_http_sdk.MQClient
}

//new实例
func newAliyunMq(diName string) queue.Queue {
	m := new(AliyunMq)
	m.client = aliyunmq.GetAliyunMq(diName)

	return m
}

//单例模式
func GetAliyunRocketQueue(diName string) queue.Queue {
	key := diName
	mu.RLock()
	q, ok := mp[key]
	mu.RUnlock()
	if ok {
		return q
	}

	q = newAliyunMq(diName)

	mu.Lock()
	mp[key] = q
	mu.Unlock()
	return q
}

/**
 * 队列消息入队
 * args[0] instanceId
 */
func (m *AliyunMq) Enqueue(ctx context.Context, key string, message string, args ...interface{}) (bool, error) {
	instanceId, _, _ := getOption(args...)

	// 获取rocketmq的producer，这个和mns不同，区分了producer和consumer，alimns统一为client
	mqProducer := m.client.GetProducer(instanceId, key)

	// aliyunmq消息格式 可以设置MessageTag和Properties等信息，先只提供最基本的MessageBody
	mqMsg := mq_http_sdk.PublishMessageRequest{
		MessageBody: message,
	}

	_, err := mqProducer.PublishMessage(mqMsg)
	if err != nil {
		return false, err
	}

	return true, nil
}

/**
* 队列消息出队
* param 第二个参数是队列名称，args[0]是instanceId，args[1]是groupId，目前只有rocketmq需要groupId
* return 第一个参数是消息 第二个参数是aliyunmq的ReceiptHandle命名为token，通过token确定消息是否从队列删除，第三个参数为消费次数
 */
func (m *AliyunMq) Dequeue(ctx context.Context, key string, args ...interface{}) (message string, token string, dequeueCount int64, err error) {
	instanceId, groupId, messageTag := getOption(args...)

	// 获取rocketmq的consumer
	mqConsumer := m.client.GetConsumer(instanceId, key, groupId, messageTag)

	//endChan := make(chan int)
	respChan := make(chan mq_http_sdk.ConsumeMessageResponse)
	errChan := make(chan error)

	go func() {
		// 长轮询消费消息
		// 长轮询表示如果topic没有消息则请求会在服务端挂住3s，3s内如果有消息可以消费则立即返回
		mqConsumer.ConsumeMessage(respChan, errChan,
			1, // 一次最多消费条数(最多可设置为16条)
			3, // 长轮询时间3秒（最多可设置为30秒）
		)
	}()

	select {
	case resp := <-respChan:
		{
			// 处理业务逻辑
			var handles []string
			respLen := len(resp.Messages)
			fmt.Printf("AliRocketMq Consume %d messages---->\n", respLen)
			if respLen != 1 {
				// 如果消息内容多于一条 可以给出提示or返回err
			}

			for _, v := range resp.Messages {
				handles = append(handles, v.ReceiptHandle)
				//fmt.Printf("\tMessageID: %s, PublishTime: %d, MessageTag: %s\n"+
				//	"\tConsumedTimes: %d, FirstConsumeTime: %d, NextConsumeTime: %d\n"+
				//	"\tBody: %s\n"+
				//	"\tProps: %s\n",
				//	v.MessageId, v.PublishTime, v.MessageTag, v.ConsumedTimes,
				//	v.FirstConsumeTime, v.NextConsumeTime, v.MessageBody, v.Properties)
				return v.MessageBody, v.ReceiptHandle, v.ConsumedTimes, nil
			}

		}
	case errMsg := <-errChan:
		{
			// 没有消息
			err = errMsg
			if strings.Contains(errMsg.(errors.ErrCode).Error(), "MessageNotExist") {
				err = nil
				// fmt.Println("\nNo new message, continue!")
			} else {
				fmt.Println("aliyunmq get msg error:", errMsg)
				time.Sleep(time.Duration(3) * time.Second)
			}
		}
	case <-time.After(35 * time.Second):
		{
			fmt.Println("Timeout of consumer message ??")
			err = errors.New("Timeout of consumer message")
		}
	}

	return
}

/**
 * 队列消息批量入队
 * args[0] instanceId
 * 注：rocket其实没有批量函数，所以循环调用publishMsg方法
 */
func (m *AliyunMq) BatchEnqueue(ctx context.Context, key string, messageList []string, args ...interface{}) (bool, error) {
	if len(messageList) == 0 {
		return false, errors.New("messageList is empty")
	}

	for _, message := range messageList {
		flag, err := m.Enqueue(ctx, key, message, args)
		if flag == false || err != nil {
			return flag, err
		}
	}

	return true, nil
}

/**
 * 确认消息接收
 * args[0]是instanceId，args[1]是groupId，args[2]是messageTag
 */
func (m *AliyunMq) AckMsg(ctx context.Context, key string, token string, args ...interface{}) (bool, error) {
	if len(token) < 1 {
		return false, errors.New("token empty")
	}

	instanceId, groupId, messageTag := getOption(args...)

	// 获取rocketmq的consumer
	mqConsumer := m.client.GetConsumer(instanceId, key, groupId, messageTag)

	var handles []string
	// rocketmq的确认函数是需要传递handle数组
	handles = append(handles, token)

	ackErr := mqConsumer.AckMessage(handles)
	if ackErr != nil {
		// 某些消息的句柄可能超时了会导致确认不成功
		fmt.Println("aliyunmq ack error token", token, ",err:", ackErr)

		for _, errAckItem := range ackErr.(errors.ErrCode).Context()["Detail"].([]mq_http_sdk.ErrAckItem) {
			fmt.Printf("aliyunmq handle ack item: \tErrorHandle:%s, ErrorCode:%s, ErrorMsg:%s\n",
				errAckItem.ErrorHandle, errAckItem.ErrorCode, errAckItem.ErrorMsg)
		}
		return false, ackErr
		//time.Sleep(time.Duration(3) * time.Second)
	} else {
		fmt.Printf("aliyunmq Ack ---->\n\t%s\n", handles)
	}

	return true, nil
}

// 缺省参数统一获取
// args[0]是instanceId，args[1]是groupId，args[2]是messageTag
func getOption(args ...interface{}) (instanceId, groupId, messageTag string) {
	instanceId = ""
	groupId = ""
	messageTag = ""

	l := len(args)
	if l > 0 {
		tempInstance, ok := args[0].(string)
		if ok {
			instanceId = tempInstance
		}
		if l > 1 {
			tempGroup, ok := args[1].(string)
			if ok {
				groupId = tempGroup
			}
		}
		if l > 2 {
			tempTag, ok := args[2].(string)
			if ok {
				messageTag = tempTag
			}
		}
	}
	return
}

func init() {
	mp = make(map[string]queue.Queue)
	queue.Register(queue.DriverTypeAliyunMq, GetAliyunRocketQueue)
}
