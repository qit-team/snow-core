package rediscache

import (
	"github.com/qit-team/snow-core/redis"
	"fmt"
	"github.com/qit-team/snow-core/config"
	"github.com/qit-team/snow-core/cache"
	"testing"
	"context"
)

var c cache.Cache

func init() {
	var err error
	redisConf := config.RedisConfig{
		Master: config.RedisBaseConfig{
			Host: "127.0.0.1",
			Port: 6379,
		},
	}

	//注册redis类
	err = redis.Pr.Register("redis", redisConf)
	if err != nil {
		fmt.Println(err)
	}

	c = cache.GetCache("redis", cache.DriverTypeRedis)
}

func TestGetSet(t *testing.T) {
	ctx := context.TODO()
	key := "cache-test"
	value := "111"
	ok, err := c.Set(ctx, key, value)
	if err != nil {
		t.Error(err)
		return
	} else if !ok {
		t.Error("set is not ok")
		return
	}

	v, err := c.Get(ctx, key)
	if err != nil {
		t.Error(err)
		return
	} else if v != value {
		t.Error("get is not same", v)
		return
	}

	ok, err = c.Delete(ctx, key)
	if err != nil {
		t.Error(err)
		return
	} else if !ok {
		t.Error("delete is not ok")
		return
	}
}
