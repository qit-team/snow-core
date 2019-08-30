package rediscache

import (
	"context"
	"fmt"
	"github.com/qit-team/snow-core/cache"
	"github.com/qit-team/snow-core/config"
	"github.com/qit-team/snow-core/redis"
	"github.com/qit-team/snow-core/utils"
	"testing"
	"time"
)

var c cache.Cache
var m *cache.BaseCache
var ctx context.Context

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

// 测试basecache
func TestBaseCache_Get_Set_IsExist(t *testing.T) {
	m := new(cache.BaseCache)
	key := "test-snow-" + fmt.Sprint(utils.GetCurrentTime())

	// key不存在情况下读取数据
	s, err := m.Get(ctx, key)
	if err != nil {
		t.Errorf("Get %s err:%s", key, err.Error())
		return
	} else if s != "" {
		t.Errorf("Get %s is not empty", key)
		return
	}

	// 判断key是否存在
	ok, err := m.IsExist(ctx, key)
	if err != nil {
		t.Errorf("IsExist %s err:%s", key, err.Error())
		return
	} else if ok {
		t.Errorf("IsExist %s is not equal false", key)
		return
	}

	value := "1"
	// 对key进行set操作且过期时间1秒
	ok, err = m.Set(ctx, key, value, 1)
	if err != nil {
		t.Errorf("Set %s err:%s", key, err.Error())
		return
	} else if !ok {
		t.Errorf("Set %s is not ok", key)
		return
	}

	// set完之后马上执行get操作
	s, _ = m.Get(ctx, key)
	if s != value {
		t.Errorf("Get %s value(%s) is not equal %s", key, s, value)
		return
	}

	time.Sleep(time.Second)

	// 一秒之后再取值，因为set时候设置过期时间为1s，如果拿不到值是正常情况
	s, _ = m.Get(ctx, key)
	if s != "" {
		t.Errorf("Get %s is not empty", key)
		return
	}
}

func TestBaseCache_Delete(t *testing.T) {
	m := new(cache.BaseCache)
	key := "test-snow1" + fmt.Sprint(utils.GetCurrentTime())
	value := "1"
	m.Set(ctx, key, value)

	ok, err := m.Delete(ctx, key)
	if err != nil {
		t.Errorf("Delete %s err:%s", key, err.Error())
		return
	} else if !ok {
		t.Errorf("Delete %s is not ok", key)
		return
	}

	s, _ := m.Get(ctx, key)
	if s != "" {
		t.Errorf("Get %s is not empty", key)
		return
	}
}

func TestBaseCache_SetMulti_GetMulti_DeleteMulti(t *testing.T) {
	m := new(cache.BaseCache)

	time := fmt.Sprint(utils.GetCurrentTime())
	key2 := "test2-snow" + time
	key3 := "test3-snow" + time
	value := "1"

	items := map[string]interface{}{
		key2: value,
		key3: value,
	}
	ok, err := m.SetMulti(ctx, items, 1)
	if err != nil {
		t.Errorf("SetMulti err:%s", err.Error())
		return
	} else if !ok {
		t.Errorf("SetMulti is not ok, keys:%s, %s", key2, key3)
		return
	}

	retMulti, err := m.GetMulti(ctx, key2, key3)
	if err != nil {
		t.Errorf("GetMulti err:%s", err.Error())
		return
	} else {
		for k, v := range retMulti {
			if v != value {
				t.Errorf("GetMulti %s value(%s) is not equal %s", k, v, value)
				return
			}
		}
	}

	ok, err = m.DeleteMulti(ctx, key2, key3)
	if err != nil {
		t.Errorf("DeleteMulti err:%s", err.Error())
		return
	} else if !ok {
		t.Errorf("DeleteMulti is not ok, keys:%s, %s", key2, key3)
		return
	}

	retMulti, err = m.GetMulti(ctx, key2, key3)
	if err != nil {
		t.Errorf("GetMulti After Delete err:%s", err.Error())
		return
	} else {
		for k, v := range retMulti {
			if v != "" {
				t.Errorf("GetMulti After Delete %s value(%s) is not empty", k, v)
				return
			}
		}
	}
}

func TestBaseCache_Expire(t *testing.T) {
	m := new(cache.BaseCache)
	key := "test-snow5-" + fmt.Sprint(utils.GetCurrentTime())

	value := "1"
	// 对key进行set操作且不设置过期时间
	ok, err := m.Set(ctx, key, value)
	if err != nil {
		t.Errorf("Set %s err:%s", key, err.Error())
		return
	} else if !ok {
		t.Errorf("Set %s is not ok", key)
		return
	}

	// 通expire函数设置过期时间
	ok, err = m.Expire(ctx, key, 1)
	if err != nil {
		t.Errorf("Expire %s err:%s", key, err.Error())
		return
	} else if !ok {
		t.Errorf("Expire %s is not ok", key)
		return
	}

	// set完之后马上执行get操作
	s, _ := m.Get(ctx, key)
	if s != value {
		t.Errorf("Get after expire %s value(%s) is not equal %s", key, s, value)
		return
	}

	time.Sleep(time.Second)

	// 一秒之后再取值，expire设置的过期时间为1s，如果拿不到值是正常情况
	s, _ = m.Get(ctx, key)
	if s != "" {
		t.Errorf("Get after expire and wait 1s %s is not empty", key)
		return
	}
}
