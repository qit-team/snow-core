package cache

import (
	"github.com/qit-team/snow-core/redis"
	"github.com/qit-team/snow-core/config"
	"fmt"
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

	Register("mock", getMockCache)
}

func getMockCache(diName string) Cache {
	return nil
}

func TestRegister(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Errorf("repeat register do not panic")
		}
	}()
	Register("mock", getMockCache)
}

func TestRegister_EmptyDriver(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Errorf("nil driver do not panic")
		}
	}()
	Register("mock", nil)
}

func TestGetCache_Empty(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Errorf("unknown driver do not panic")
		}
	}()
	GetCache("redis", "empty")
}

func TestGetCache_Nil(t *testing.T) {
	defer func() {
		if e := recover(); e == nil {
			t.Errorf("unknown diName do not panic")
		}
	}()
	GetCache("unknown", "mock")
}
