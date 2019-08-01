package rediscache

import (
	"github.com/qit-team/snow-core/redis"
	"fmt"
	"github.com/qit-team/snow-core/config"
	"github.com/qit-team/snow-core/cache"
	"testing"
	"context"
	"time"
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

func TestGetSetDelete(t *testing.T) {
	c := cache.GetCache("redis", cache.DriverTypeRedis)
	ctx := context.TODO()
	key := "test-cache"
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

	v, err = c.Get(ctx, key)
	if err != nil {
		t.Error(err)
		return
	} else if v != "" {
		t.Errorf("delete %s failed", key)
		return
	}
}

func TestSetMultiAndGetMulti(t *testing.T) {
	ctx := context.TODO()
	items := map[string]interface{}{
		"test-key1": "111",
		"test-key2": "222",
	}
	_, err := c.SetMulti(ctx, items, 1)
	if err != nil {
		t.Error(err)
		return
	}

	m, err := c.GetMulti(ctx, "test-key1", "test-key2")
	if err != nil {
		t.Error(err)
		return
	} else if len(m) != 2 {
		t.Error("get values's length is not enough")
		return
	}
	var value interface{}
	var ok bool
	for k, v := range m {
		if value, ok = items[k]; !ok {
			t.Errorf("key %s is not exist", k)
			return
		}
		if value != v {
			t.Errorf("key %s is not same", k)
			return
		}
	}

	time.Sleep(time.Millisecond * 1100)
	m, err = c.GetMulti(ctx, "test-key1", "test-key2")
	if err != nil {
		t.Error(err)
		return
	} else if len(m) != 2 {
		t.Error("get values's length is not enough")
		return
	}

	for k, v := range m {
		if _, ok = items[k]; !ok {
			t.Errorf("key %s is not exist", k)
			return
		}
		if v != "" {
			t.Errorf("key %s is not empty", k)
			return
		}
	}
}

func TestDeleteMulti(t *testing.T) {
	ctx := context.TODO()
	items := map[string]interface{}{
		"test-key3": "111",
		"test-key4": "222",
	}

	c.SetMulti(ctx, items)

	_, err := c.DeleteMulti(ctx, "test-key3", "test-key4")
	if err != nil {
		t.Error(err)
		return
	}

	var ok bool
	m, err := c.GetMulti(ctx, "test-key3", "test-key4")
	if err != nil {
		t.Error(err)
		return
	} else if len(m) != 2 {
		t.Error("get values's length is not enough")
		return
	}

	for k, v := range m {
		if _, ok = items[k]; !ok {
			t.Errorf("key %s is not exist", k)
			return
		}
		if v != "" {
			t.Errorf("key %s is not empty", k)
			return
		}
	}
}

func TestExpireExist(t *testing.T) {
	ctx := context.TODO()
	key := "test-expire"
	value := "222"
	c.Set(ctx, key, value)

	ok, err := c.IsExist(ctx, key)
	if err != nil {
		t.Error(err)
		return
	} else if !ok {
		t.Errorf("key %s is not exist", key)
		return
	}

	c.Expire(ctx, key, 1)
	time.Sleep(time.Millisecond * 1100)

	ok, err = c.IsExist(ctx, key)
	if err != nil {
		t.Error(err)
		return
	} else if ok {
		t.Errorf("key %s is exist", key)
		return
	}
}
