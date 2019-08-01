package queue

import (
	"github.com/qit-team/snow-core/config"
	"fmt"
	"github.com/qit-team/snow-core/redis"
	"testing"
)

func init() {
	redisConf := config.RedisConfig{
		Master: config.RedisBaseConfig{
			Host: "127.0.0.1",
			Port: 6379,
		},
	}

	//注册redis类
	err := redis.Pr.Register("redis", redisConf)
	if err != nil {
		fmt.Println(err)
	}

	Register("mock", getMockQueue)
}

func getMockQueue(diName string) Queue {
	return nil
}

func TestRegister(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Errorf("repeat register do not panic")
		}
	}()
	Register("mock", getMockQueue)
}

func TestRegister_EmptyDriver(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Errorf("nil driver do not panic")
		}
	}()
	Register("mock", nil)
}

func TestGetQueue_Empty(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Errorf("unknown driver do not panic")
		}
	}()
	GetQueue("redis", "empty")
}

func TestGetQueue_Nil(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Errorf("unknown diName do not panic")
		}
	}()
	GetQueue("unknown", "mock")
}
