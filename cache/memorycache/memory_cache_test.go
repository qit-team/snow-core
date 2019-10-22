package memorycache

import (
	"context"
	"github.com/qit-team/snow-core/cache"
	"testing"
	"time"
	"errors"
)

var c cache.Cache

func init() {
	c = cache.GetCache("memory", cache.DriverTypeMemory)
}

func TestGetSetDelete(t *testing.T) {
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

func TestMemoryCache_IncrBy(t *testing.T) {
	ctx := context.TODO()
	key := "test-incr"
	value := "ab"
	c.Set(ctx, key, value)

	_, err := c.IncrBy(ctx, key, 3)
	if !errors.Is(err, ErrWrongDataType) {
		t.Errorf("wrong error type %s", err)
		return
	}

	c.Set(ctx, key, 400)
	res, err := c.IncrBy(ctx, key, 3)
	if err != nil {
		t.Error(err)
		return
	} else if res != 403 {
		t.Errorf("wrong increment %d", res)
		return
	}

	c.Delete(ctx, key)
	res, err = c.IncrBy(ctx, key, -30)
	if err != nil {
		t.Error(err)
		return
	} else if res != -30 {
		t.Errorf("wrong increment %d", res)
		return
	}
}


func TestMemoryCache_DecrBy(t *testing.T) {
	ctx := context.TODO()
	key := "test-desc"
	value := "ab"
	c.Set(ctx, key, value)

	_, err := c.DecrBy(ctx, key, 10)
	if !errors.Is(err, ErrWrongDataType) {
		t.Errorf("wrong error type %s", err)
		return
	}

	c.Set(ctx, key, 400)
	res, err := c.DecrBy(ctx, key, 10)
	if err != nil {
		t.Error(err)
		return
	} else if res != 390 {
		t.Errorf("wrong decrement %d", res)
		return
	}

	c.Delete(ctx , key)
	res, err = c.DecrBy(ctx, key, -30)
	if err != nil {
		t.Error(err)
		return
	} else if res != 30 {
		t.Errorf("wrong decrement %d", res)
		return
	}
}



