package rocketqueue

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/gogap/errors"
	"github.com/qit-team/snow-core/log/logger"
	"github.com/qit-team/snow-core/queue"
	rkmq "github.com/qit-team/snow-core/rocketmq"
)

var (
	mp map[string]queue.Queue
	mu sync.RWMutex
)

type RocketQueue struct {
	Consumer rocketmq.PushConsumer
	Producer rocketmq.Producer

	consumerMessageChan chan *primitive.MessageExt
	//producerMessageChan chan *primitive.MessageExt

	consumerOnce sync.Once
	producerOnce sync.Once
}

func (m *RocketQueue) initProducer(ctx context.Context) error {
	var err error
	m.producerOnce.Do(
		func() {
			err = m.Producer.Start()
			if err != nil {
				logger.Fatal(ctx, "RocketQueue:Producer:Start", err.Error())
				return
			}
		})
	return err
}

func (m *RocketQueue) initConsumer(ctx context.Context, topic, messageTag string, num int) error {
	var err error
	m.consumerOnce.Do(
		func() {
			m.consumerMessageChan = make(chan *primitive.MessageExt, num)

			var selector consumer.MessageSelector
			if len(messageTag) > 0 {
				selector = consumer.MessageSelector{
					Type:       consumer.TAG,
					Expression: messageTag,
				}
			}
			err = m.Consumer.Subscribe(topic, selector, func(ctx context.Context, messages ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
				// 取到的消息放入管道，交给下游处理
				for _, msg := range messages {
					m.consumerMessageChan <- msg
				}

				return consumer.ConsumeSuccess, nil
			})
			if err != nil {
				logger.Error(ctx, "RocketQueue:Subscribe", err.Error())
				return
			}

			err = m.Consumer.Start()
			if err != nil {
				logger.Fatal(ctx, "RocketQueue:Start", err.Error())
				return
			}

			go func() {
				var sigs = []os.Signal{
					syscall.SIGHUP,
					syscall.SIGUSR1,
					syscall.SIGUSR2,
					syscall.SIGINT,
					syscall.SIGTERM,
					syscall.SIGQUIT,
				}
				c := make(chan os.Signal)
				signal.Notify(c, sigs...)
				for {
					sig := <-c //blocked
					switch sig {
					case syscall.SIGINT, syscall.SIGTERM:
						close(m.consumerMessageChan)
						err = m.Consumer.Shutdown()
						if err != nil {
							logger.Error(ctx, "Shutdown.Failure", err.Error())
							return
						}
						return
					default:
					}
				}
				fmt.Println("停止订阅消息")
			}()
		})
	if err != nil {
		logger.Error(ctx, "RocketQueue:initConsumer", err.Error())
		return err
	}

	return nil
}

// new实例
func newRocketQueue(diName string) queue.Queue {
	m := new(RocketQueue)
	client := rkmq.GetRocketMq(diName)

	m.Producer = client.Producer
	m.Consumer = client.Consumer

	return m
}

// GetRocketQueue
//
// 单例模式
func GetRocketQueue(diName string) queue.Queue {
	key := diName
	mu.RLock()
	q, ok := mp[key]
	mu.RUnlock()
	if ok {
		return q
	}

	q = newRocketQueue(diName)

	mu.Lock()
	mp[key] = q
	mu.Unlock()
	return q
}

// Enqueue 队列消息入队
//
// args[0] instanceId
func (m *RocketQueue) Enqueue(ctx context.Context, key string, message string, args ...interface{}) (bool, error) {
	err := m.initProducer(ctx)
	if err != nil {
		return false, err
	}
	_, _, messageTag, timeLevel := getOption(args...)
	log.Printf("messageTag: %v", messageTag)
	if len(messageTag) > 0 {
		tags := strings.Split(messageTag, "||")
		for i := 0; i < len(tags); i++ {
			tag := strings.Trim(tags[i%3], " ")
			msg := &primitive.Message{
				Topic: key,
				Body:  []byte(message),
			}
			msg.WithTag(tag)
			// https://rocketmq.apache.org/docs/4.x/producer/04message3/
			if timeLevel > 0 && timeLevel <= 18 {
				msg.WithDelayTimeLevel(timeLevel)
			}
			log.Printf("send for tag: %v", tag)
			res, err := m.Producer.SendSync(context.Background(), msg)
			if err != nil {
				return false, err
			}
			//logger.Info(ctx, "Enqueue", res.String())
			log.Printf("Enqueue: %s %v", message, res.MsgID)
		}
	} else {
		msg := &primitive.Message{
			Topic: key,
			Body:  []byte(message),
		}
		// https://rocketmq.apache.org/docs/4.x/producer/04message3/
		if timeLevel > 0 && timeLevel <= 18 {
			msg.WithDelayTimeLevel(timeLevel)
		}
		res, err := m.Producer.SendSync(ctx, msg)
		if err != nil {
			return false, err
		}
		//logger.Info(ctx, "Enqueue", res.String())
		log.Printf("Enqueue: %s %v", message, res.MsgID)
	}

	return true, nil
}

// Dequeue 队列消息出队
//
// param 第二个参数是队列名称，args[0]是instanceId，args[1]是groupId，目前只有rocketmq需要groupId
//
// return 第一个参数是消息 第二个参数是aliyunmq的ReceiptHandle命名为token，通过token确定消息是否从队列删除，第三个参数为消费次数
func (m *RocketQueue) Dequeue(ctx context.Context, key string, args ...interface{}) (message string, tag string, token string, dequeueCount int64, err error) {
	_, _, messageTag, _ := getOption(args...)

	err = m.initConsumer(ctx, key, messageTag, 5)
	if err != nil {
		return
	}

	select {
	case msg, ok := <-m.consumerMessageChan:
		if !ok {
			return "", "", "", 0, nil
		}
		return string(msg.Body), msg.GetTags(), "", int64(msg.ReconsumeTimes), nil
	}
}

// BatchEnqueue 队列消息批量入队
// args[0] instanceId
// 注：rocket其实没有批量函数，所以循环调用publishMsg方法
func (m *RocketQueue) BatchEnqueue(ctx context.Context, key string, messageList []string, args ...interface{}) (bool, error) {
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

// AckMsg 确认消息接收
// args[0]是instanceId，args[1]是groupId，args[2]是messageTag, args[2]是delayTimeLevel
func (m *RocketQueue) AckMsg(ctx context.Context, key string, token string, args ...interface{}) (bool, error) {
	return true, nil
}

// getOption 缺省参数统一获取
//
// args[0]是instanceId，args[1]是groupId，args[2]是messageTag, args[3]是delayTimeLevel
func getOption(args ...interface{}) (instanceId, groupId, messageTag string, delayTimeLevel int) {
	instanceId = ""
	groupId = ""
	messageTag = ""
	delayTimeLevel = 0
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
		if l > 3 {
			tempDelayTimeLevel, ok := args[3].(int)
			if ok {
				delayTimeLevel = tempDelayTimeLevel
			}
		}
	}
	return
}

func init() {
	mp = make(map[string]queue.Queue)
	queue.Register(queue.DriverTypeRocketMq, GetRocketQueue)
}
