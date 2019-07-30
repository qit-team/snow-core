package redisqueue

import (
	"testing"
	"context"
	"github.com/qit-team/snow-core/config"
	"fmt"
	"github.com/qit-team/snow-core/redis"
	"github.com/qit-team/snow-core/queue"
)

var q queue.Queue

func init() {
	redisConf := config.RedisConfig{
		Master: config.RedisBaseConfig{
			Host: "127.0.0.1",
			Port: 6379,
		},
	}

	//注册redis类
	err := redis.Pr.Register("redis", redisConf, true)
	if err != nil {
		fmt.Println(err)
	}

	q = queue.GetQueue("redis", queue.DriverTypeRedis)
}

func TestEnqueue(t *testing.T) {
	q := queue.GetQueue("redis", queue.DriverTypeRedis)
	topic := "snow-topic-one"
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

	message, _, err = q.Dequeue(ctx, topic)
	if err != nil {
		t.Error(err)
		return
	} else if message != "" {
		t.Error("message must be empty")
		return
	}
}

func TestBatchEnqueue(t *testing.T) {
	ctx := context.TODO()
	topic := "snow-topic-batch"
	messages := []string{"11", "21"}
	_, err := q.BatchEnqueue(ctx, topic, messages)
	if err != nil {
		t.Error("batch enqueue error", err)
		return
	}

	fmt.Println("batch enqueue", topic, messages)

	message, _, err := q.Dequeue(ctx, topic)
	if err != nil {
		t.Error(err)
		return
	}
	if message != messages[0] {
		t.Errorf("message is not same origin:%s real:%s", messages[0], message)
		return
	}

	message, _, err = q.Dequeue(ctx, topic)
	if err != nil {
		t.Error(err)
		return
	}
	if message != messages[1] {
		t.Errorf("message is not same origin:%s real:%s", messages[1], message)
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
