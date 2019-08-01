package alimnsqueue

import (
	"fmt"
	"github.com/qit-team/snow-core/queue"
	"github.com/qit-team/snow-core/config"
	"testing"
	"context"
	"github.com/qit-team/snow-core/alimns"
	"github.com/qit-team/snow-core/utils"
	"io/ioutil"
	"strings"
)

var q queue.Queue

func init() {
	//需要自己在文件填好配置
	bs, err := ioutil.ReadFile("../../.env.mns")
	conf := config.MnsConfig{}
	if err == nil {
		str := string(bs)
		arr := strings.Split(str, "\n")
		if len(arr) >= 3 {
			conf.Url = arr[0]
			conf.AccessKeyId = arr[1]
			conf.AccessKeySecret = arr[2]
		}
	}

	//注册alimns类
	err = alimns.Pr.Register("ali_mns", conf)
	if err != nil {
		fmt.Println(err)
	}

	q = queue.GetQueue("ali_mns", queue.DriverTypeAliMns)
}

func TestEnqueue(t *testing.T) {
	q := queue.GetQueue("ali_mns", queue.DriverTypeAliMns)
	topic := "snow-topic-one" + fmt.Sprint(utils.GetCurrentTime())
	ctx := context.TODO()
	msg := "1"
	ok, err := q.Enqueue(ctx, topic, msg)
	if err != nil {
		t.Error(err)
		return
	}
	if !ok {
		t.Error("enqueue is not ok")
		return
	}

	message, token, err := q.Dequeue(ctx, topic)
	if err != nil {
		t.Error(err)
		return
	}
	if message != msg {
		t.Errorf("message is not same %s", message)
		return
	}

	ok, err = q.AckMsg(ctx, topic, token)
	if err != nil {
		t.Error(err)
		return
	}
	if !ok {
		t.Error("ack is not ok")
		return
	}

	message, token, err = q.Dequeue(ctx, topic)
	if err != nil {
		t.Error(err)
		return
	} else if message != "" {
		t.Error("message from blank queue must be empty")
		return
	}

	_, err = q.AckMsg(ctx, topic, token)
	if !(err != nil && err.Error() == "token empty") {
		t.Error("must return empty ack token error")
	}
}

func TestBatchEnqueue(t *testing.T) {
	ctx := context.TODO()
	topic := "snow-topic-batch" + fmt.Sprint(utils.GetCurrentTime())
	messages := []string{"11", "21"}
	_, err := q.BatchEnqueue(ctx, topic, messages)
	if err != nil {
		t.Error("batch enqueue error", err)
		return
	}

	fmt.Println("batch enqueue", topic, messages)

	message1, token1, err := q.Dequeue(ctx, topic)
	if err != nil {
		t.Error(err)
		return
	}

	message2, token2, err := q.Dequeue(ctx, topic)
	if err != nil {
		t.Error(err)
		return
	}

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

	ok, err := q.AckMsg(ctx, topic, token1)
	if err != nil {
		t.Errorf("message1 ack err:%s", err.Error())
		return
	}
	if !ok {
		t.Error("message1 ack is not ok")
		return
	}

	ok, err = q.AckMsg(ctx, topic, token2)
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
	topic := "snow-topic-batch"
	messages := make([]string, 0)
	_, err := q.BatchEnqueue(ctx, topic, messages)
	if err == nil {
		t.Error("empty message must return error")
		return
	}
}

func Test_getOption(t *testing.T) {
	delay, priority := getOption(int64(1), int64(10))
	if delay != 1 {
		t.Errorf("delay is not equal 1. %d", delay)
	} else if priority != 10 {
		t.Errorf("priority is not equal 10. %d", priority)
	}
}
