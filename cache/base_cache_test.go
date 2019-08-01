package cache

import (
	"testing"
	"context"
	"github.com/qit-team/snow-core/redis"
	"github.com/qit-team/snow-core/config"
	"fmt"
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

//func TestBaseCache_Get_IsExist_Set(t *testing.T) {
//	key := "test-" + fmt.Sprint(utils.GetCurrentTime())
//	s, err := m.Get(ctx, key)
//	if err != nil {
//		t.Errorf("Get %s err:%s", key, err.Error())
//		return
//	} else if s != "" {
//		t.Errorf("Get %s is not empty", key)
//		return
//	}
//
//	ok, err := m.IsExist(ctx, key)
//	if err != nil {
//		t.Errorf("IsExist %s err:%s", key, err.Error())
//		return
//	} else if ok {
//		t.Errorf("IsExist %s is not equal false", key)
//		return
//	}
//
//	value := "1"
//	ok, err = m.Set(ctx, key, value, 1)
//	if err != nil {
//		t.Errorf("Set %s err:%s", key, err.Error())
//		return
//	} else if !ok {
//		t.Errorf("Set %s is not ok", key)
//		return
//	}
//
//	s, _ = m.Get(ctx, key)
//	if s != value {
//		t.Errorf("Get %s value(%s) is not equal %s", key, s, value)
//		return
//	}
//
//	time.Sleep(time.Second)
//
//	s, _ = m.Get(ctx, key)
//	if s != "" {
//		t.Errorf("Get %s is not empty", key)
//		return
//	}
//}
//
//func TestBaseCache_Delete(t *testing.T) {
//	key := "test1-" + fmt.Sprint(utils.GetCurrentTime())
//	value := "1"
//	m.Set(ctx, key, value)
//
//	ok, err := m.Delete(ctx, key)
//	if err != nil {
//		t.Errorf("Delete %s err:%s", key, err.Error())
//		return
//	} else if !ok {
//		t.Errorf("Delete %s is not ok", key)
//		return
//	}
//
//	s, _ := m.Get(ctx, key)
//	if s != "" {
//		t.Errorf("Get %s is not empty", key)
//		return
//	}
//}
//
//func TestBaseCache_SetMulti_GetMulti_DeleteMulti(t *testing.T) {
//	time := fmt.Sprint(utils.GetCurrentTime())
//	key2 := "test2-" + time
//	key3 := "test3-" + time
//	value := "1"
//
//	items := map[string]interface{}{
//		key2: value,
//		key3: value,
//	}
//	m.SetMulti(ctx, items, 1)
//}
