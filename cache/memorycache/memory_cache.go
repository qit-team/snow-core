package memorycache

import (
	"context"
	"github.com/qit-team/snow-core/cache"
	"sync"
	"time"
	"errors"
	"strconv"
	"fmt"
)

const (
	MaxPersistenceTime = 86400 * 365 * 10
)

var (
	mp map[string]cache.Cache
	mu sync.RWMutex
	ErrWrongDataType error
)

type Item struct {
	data     interface{}
	expireAt time.Time
}

type MemoryCache struct {
	items map[string]Item
	mu sync.RWMutex
}

func init()  {
	ErrWrongDataType = errors.New("wrong data type")
}

//实例模式
func newMemoryCache() cache.Cache {
	m := new(MemoryCache)
	m.items = make(map[string]Item)
	return m
}

//单例模式
func GetRedisCache(diName string) cache.Cache {
	key := diName
	mu.RLock()
	q, ok := mp[key]
	mu.RUnlock()
	if ok {
		return q
	}

	q = newMemoryCache()
	mu.Lock()
	mp[key] = q
	mu.Unlock()
	return q
}

/**
 * 获取缓存key的数据
 * 注意事项，如果key值不存在的话，返回的是空字符串，而不是nil
 */
func (c *MemoryCache) Get(ctx context.Context, key string) (interface{}, error) {
	c.mu.RLock()
	value, ok := c.items[key]
	c.mu.RUnlock()
	if !ok {
		return "", nil
	}
	if inExpire(value.expireAt) {
		return value.data, nil
	}
	return "", nil
}

func (c *MemoryCache) GetMulti(ctx context.Context, keys ...string) (map[string]interface{}, error) {
	arr := make(map[string]interface{})
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, key := range keys {
		if value, ok := c.items[key]; ok && inExpire(value.expireAt) {
			arr[key] = value.data
		} else {
			arr[key] = ""
		}
	}
	return arr, nil
}

func (c *MemoryCache) Set(ctx context.Context, key string, value interface{}, ttl ...int) (bool, error) {
	t := cache.GetTTLOrDefault(ttl...)
	if t == 0 {
		t = MaxPersistenceTime
	}
	item := Item{
		data:     value,
		expireAt: time.Now().Add(time.Duration(t) * time.Second),
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = item
	return true, nil
}

func (c *MemoryCache) SetMulti(ctx context.Context, items map[string]interface{}, ttl ...int) (bool, error) {
	t := cache.GetTTLOrDefault(ttl...)
	if t == 0 {
		t = MaxPersistenceTime
	}
	expireAt := time.Now().Add(time.Duration(t) * time.Second)
	var item Item
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, value := range items {
		item = Item{
			data:     value,
			expireAt: expireAt,
		}
		c.items[key] = item
	}
	return true, nil
}

func (c *MemoryCache) Delete(ctx context.Context, key string) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.items[key]; ok {
		delete(c.items, key)
	}
	return true, nil
}

func (c *MemoryCache) DeleteMulti(ctx context.Context, keys ...string) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, key := range keys {
		if _, ok := c.items[key]; ok {
			delete(c.items, key)
		}
	}
	return true, nil
}

func (c *MemoryCache) Expire(ctx context.Context, key string, ttl ...int) (bool, error) {
	t := cache.GetTTLOrDefault(ttl...)
	expireAt := time.Now().Add(time.Duration(t))
	c.mu.Lock()
	defer c.mu.Unlock()
	if item, ok := c.items[key]; ok {
		if inExpire(item.expireAt) {
			item.expireAt = expireAt
			c.items[key] = item
		} else {
			delete(c.items, key)
		}
	}
	return true, nil
}

func (c *MemoryCache) IsExist(ctx context.Context, key string) (bool, error) {
	c.mu.RLock()
	value, ok := c.items[key]
	c.mu.RUnlock()
	if ok && inExpire(value.expireAt) {
		return true, nil
	}
	return false, nil
}

func (c *MemoryCache) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var newValue int64
	if item, ok := c.items[key]; ok {
		if val, err := interfaceToInt64(item.data); err == nil {
			newValue = val + value
			item.data = newValue
			c.items[key] = item
		} else {
			return 0, ErrWrongDataType
		}
	} else {
		newValue = value
		item = Item{
			data:     newValue,
			expireAt: time.Now().Add(time.Duration(MaxPersistenceTime) * time.Second),
		}
		c.items[key] = item
	}
	return newValue, nil
}

func (c *MemoryCache) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var newValue int64
	if item, ok := c.items[key]; ok {
		if val, err := interfaceToInt64(item.data); err == nil {
			newValue = val - value
			item.data = newValue
			c.items[key] = item
		} else {
			return 0, ErrWrongDataType
		}
	} else {
		newValue = -value
		item = Item{
			data:     newValue,
			expireAt: time.Now().Add(time.Duration(MaxPersistenceTime) * time.Second),
		}
		c.items[key] = item
	}
	return newValue, nil
}

func inExpire(u time.Time) bool {
	return time.Now().Before(u)
}

func interfaceToInt64(value interface{}) (int64, error) {
	v := fmt.Sprintf("%d", value)
	val, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}
	return int64(val), nil
}

func init() {
	mp = make(map[string]cache.Cache)
	cache.Register(cache.DriverTypeMemory, GetRedisCache)
}
