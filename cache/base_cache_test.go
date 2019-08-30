package cache

import (
	"context"
	"fmt"
	"github.com/qit-team/snow-core/config"
	"github.com/qit-team/snow-core/redis"
	"testing"
)

var m *BaseCache
var ctx context.Context

func init() {
	m = new(BaseCache)
	m.Prefix = "test:"
	ctx = context.TODO()

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
}

func TestBaseCache_GetPrefixOrDefault(t *testing.T) {
	m := new(BaseCache)
	s1 := m.GetPrefixOrDefault()
	if s1 != DefaultPrefix {
		t.Errorf("GetPrefixOrDefault is not equal default:%s", DefaultPrefix)
		return
	}

	m.Prefix = "m:"
	s2 := m.GetPrefixOrDefault()
	if s2 != m.Prefix {
		t.Errorf("GetPrefixOrDefault is not equal default:%s", m.Prefix)
		return
	}
}

func TestBaseCache_GetDiNameOrDefault(t *testing.T) {
	m := new(BaseCache)
	s1 := m.GetDiNameOrDefault()
	if s1 != DefaultDiName {
		t.Errorf("GetDiNameOrDefault is not equal default:%s", DefaultDiName)
		return
	}

	m.DiName = "di"
	s2 := m.GetDiNameOrDefault()
	if s2 != m.DiName {
		t.Errorf("GetDiNameOrDefault is not equal %s", m.DiName)
		return
	}
}

func TestBaseCache_GetDriverTypeOrDefault(t *testing.T) {
	m := new(BaseCache)
	s1 := m.GetDriverTypeOrDefault()
	if s1 != DefaultDriverType {
		t.Errorf("GetDriverTypeOrDefault is not equal default:%s", DefaultDriverType)
		return
	}

	m.DriverType = "dr"
	s2 := m.GetDriverTypeOrDefault()
	if s2 != m.DriverType {
		t.Errorf("GetDriverTypeOrDefault is not equal %s", m.DriverType)
		return
	}
}

func TestBaseCache_GetTTLOrDefault(t *testing.T) {
	m := new(BaseCache)
	t1 := m.GetTTLOrDefault()
	if t1 != DefaultTTL {
		t.Errorf("GetTTLOrDefault is not equal default:%d", DefaultTTL)
		return
	}

	m.SetTTL(1)
	t2 := m.GetTTLOrDefault()
	if t2 != 1 {
		t.Error("GetTTLOrDefault is not equal 1")
		return
	}

	m.SetTTL(0)
	t3 := m.GetTTLOrDefault()
	if t3 != 0 {
		t.Error("GetTTLOrDefault is not equal 0")
		return
	}
}

func TestBaseCache_getTTL(t *testing.T) {
	m := new(BaseCache)
	t1 := m.getTTL(1)
	if t1 != 1 {
		t.Error("getTTL is not equal 1")
		return
	}

	t2 := m.getTTL()
	if t2 != DefaultTTL {
		t.Errorf("getTTL is not equal %d", DefaultTTL)
		return
	}
}

func TestBaseCache_KeyRelated(t *testing.T) {
	m := new(BaseCache)
	key := "snow-test"
	redisKey := m.key(key)

	if len(redisKey) < 9 {
		t.Error("get redis key error")
		return
	}

	redisKeyList := m.keys("test-key-A", "test-key-B")

	for k, v := range redisKeyList {
		if len(v) < 10 {
			t.Error("get redis key error")
			return
		}
		fmt.Println("get redis key list:", k, v)
	}

	tempKey := m.removePrefix(redisKey)
	if len(tempKey) < 9 {
		t.Error("remove key error")
		return
	}
}
