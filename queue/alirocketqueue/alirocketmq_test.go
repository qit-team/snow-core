package alirocketqueue

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/qit-team/snow-core/aliyunmq"
	"github.com/qit-team/snow-core/config"
	"github.com/qit-team/snow-core/queue"
)

var q queue.Queue

func init() {
	// 需要自己在文件填好配置
	conf := config.AliyunMqConfig{}
	conf = getConfig()
	//注册alimns类
	err := aliyunmq.Pr.Register("aliyun_mq", conf)

	if err != nil {
		fmt.Println(err)
	}

	q = queue.GetQueue("aliyun_mq", queue.DriverTypeAliyunMq)
}

func TestEnqueue(t *testing.T) {
	q := queue.GetQueue("aliyun_mq", queue.DriverTypeAliyunMq)
	topic := "SNOW-TOPIC-TEST"
	groupId := "GID-SNOW-TOPIC-TEST"
	ctx := context.TODO()
	msg := "msg from snow core"
	ok, err := q.Enqueue(ctx, topic, msg)
	if err != nil {
		t.Error(err)
		return
	}
	if !ok {
		t.Error("enqueue is not ok")
		return
	}

	message, tag, token, dequeueCount, err := q.Dequeue(ctx, topic, "", groupId)
	fmt.Println("message content:", message)
	fmt.Println("message tag:", tag)
	fmt.Println("message dequeue num:", dequeueCount)
	fmt.Println("message token:", token)
	if err != nil {
		t.Error(err)
		return
	}
	if message != msg {
		t.Errorf("message is not same %s", message)
		return
	}

	ok, err = q.AckMsg(ctx, topic, token, "", groupId)

	fmt.Println("info:", ok, err)
	if err != nil {
		t.Error(err)
		return
	}
	if !ok {
		t.Error("ack is not ok")
		return
	}

	message, tag, token, dequeueCount, err = q.Dequeue(ctx, topic, "", groupId)
	fmt.Println("message content:", message)
	fmt.Println("message tag:", tag)
	fmt.Println("message dequeue num:", dequeueCount)
	fmt.Println("message token:", token)
	if err != nil {
		t.Error(err)
		return
	} else if message != "" {
		t.Error("message from blank queue must be empty")
		return
	}

	_, err = q.AckMsg(ctx, topic, token, "", groupId)
	fmt.Println("ackMsg,errInfo", err)
	if !(err != nil && err.Error() == "token empty") {
		t.Error("must return empty ack token error")
	}
}

func TestBatchEnqueue(t *testing.T) {
	ctx := context.TODO()
	topic := "SNOW-TOPIC-TEST"
	groupId := "GID-SNOW-TOPIC-TEST"
	messages := []string{"11", "21"}
	_, err := q.BatchEnqueue(ctx, topic, messages)
	if err != nil {
		t.Error("batch enqueue error", err)
		return
	}

	fmt.Println("batch enqueue", topic, messages)

	message1, _, token1, dequeueCount, err := q.Dequeue(ctx, topic, "", groupId)
	if err != nil {
		t.Error(err)
		return
	}

	message2, _, token2, dequeueCount, err := q.Dequeue(ctx, topic, "", groupId)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("TestBatchEnqueue:dequeueCount:", dequeueCount)

	if message1 == messages[0] {
		if message2 != messages[1] {
			t.Errorf("message2 is not same origin:%s real:%s", messages[1], message2)
			return
		}
	} else if message2 == messages[0] {
		if message1 != messages[1] {
			t.Errorf("message1 is not same origin:%s real:%s", messages[1], message1)
			return
		}
	} else {
		t.Errorf("message is not same %s", messages[1])
		return
	}

	ok, err := q.AckMsg(ctx, topic, token1, "", groupId)
	if err != nil {
		t.Errorf("message1 ack err:%s", err.Error())
		return
	}
	if !ok {
		t.Error("message1 ack is not ok")
		return
	}

	ok, err = q.AckMsg(ctx, topic, token2, "", groupId)
	if err != nil {
		t.Errorf("message1 ack err:%s", err.Error())
		return
	}
	if !ok {
		t.Error("message2 ack is not ok")
		return
	}
}

func TestBatchEnqueueEmpty(t *testing.T) {
	ctx := context.TODO()
	topic := "SNOW-TOPIC-TEST"
	groupId := "GID-SNOW-TOPIC-TEST"
	messages := make([]string, 0)
	_, err := q.BatchEnqueue(ctx, topic, messages, "", groupId)
	fmt.Println("TestBatchEnqueueEmpty.Error", err)
	if err == nil {
		t.Error("empty message must return error")
		return
	}
}

func Test_getOption(t *testing.T) {
	instanceId, groupId, _ := getOption("", "GID-SNOW-TOPIC-TEST")
	if instanceId != "" {
		t.Errorf("delay is not equal 1. %s", instanceId)
	} else if groupId != "GID-SNOW-TOPIC-TEST" {
		t.Errorf("priority is not equal 10. %s", groupId)
	}
}

func getConfig() config.AliyunMqConfig {
	//需要自己在文件填好配置
	bs, err := ioutil.ReadFile("../../.env.aliyunmq")

	conf := config.AliyunMqConfig{}
	if err == nil {
		str := string(bs)
		arr := strings.Split(str, "\n")
		if len(arr) >= 3 {
			conf.EndPoint = arr[0]
			conf.AccessKey = arr[1]
			conf.SecretKey = arr[2]
		}
	}
	return conf
}
