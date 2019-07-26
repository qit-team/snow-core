package rediscache

import (
	"context"
	redis_pool "github.com/hetiansu5/go-redis-pool"
	"github.com/qit-team/snow-core/redis"
	"github.com/qit-team/snow-core/cache"
	"sync"
)

var (
	mp map[string]cache.Cache
	mu sync.RWMutex
)

type RedisCache struct {
	client *redis_pool.ReplicaPool
}

//实例模式
func newRedisCache(diName string) cache.Cache {
	m := new(RedisCache)
	m.client = redis.GetRedis(diName)
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

	q = newRedisCache(diName)
	mu.Lock()
	mp[key] = q
	mu.Unlock()
	return q
}

/**
 * 获取缓存key的数据
 * 注意事项，如果key值不存在的话，返回的是空字符串，而不是nil
 */
func (c *RedisCache) Get(ctx context.Context, key string) (interface{}, error) {
	value, err := c.client.Get(key)
	if err == redis_pool.ErrNil {
		return "", nil
	}
	return value, err
}

func (c *RedisCache) GetMulti(ctx context.Context, keys ...string) (map[string]interface{}, error) {
	cKeys := convert(keys)
	values, err := c.client.MGet(cKeys...)
	if err != nil {
		return nil, err
	}

	arr := make(map[string]interface{})
	for index, key := range keys {
		arr[key] = values[index]
	}
	return arr, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl ...int) (bool, error) {
	t := cache.GetTTLOrDefault(ttl...)
	return c.client.SetEX(key, value, int64(t))
}

func (c *RedisCache) SetMulti(ctx context.Context, items map[string]interface{}, ttl ...int) (bool, error) {
	arr := make([]interface{}, 0)
	for key, value := range items {
		arr = append(arr, key, value)
	}
	ok, err := c.client.MSet(arr...)
	if err != nil {
		return ok, err
	}

	t := cache.GetTTLOrDefault(ttl...)
	if t > 0 {
		t64 := int64(t)
		for key, _ := range items {
			c.client.Expire(key, t64)
		}
	}
	return true, nil
}

func (c *RedisCache) Delete(ctx context.Context, key string) (bool, error) {
	res, err := c.client.Del(key)
	return res > 0, err
}

func (c *RedisCache) DeleteMulti(ctx context.Context, keys ...string) (bool, error) {
	cKeys := convert(keys)
	res, err := c.client.Del(cKeys...)
	return res > 0, err
}

func (c *RedisCache) Expire(ctx context.Context, key string, ttl ...int) (bool, error) {
	t := cache.GetTTLOrDefault(ttl...)
	return c.client.Expire(key, int64(t))
}

func (c *RedisCache) IsExist(ctx context.Context, key string) (bool, error) {
	num, err := c.client.Exists(key)
	return num == 1, err
}

func convert(keys []string) []interface{} {
	arr := make([]interface{}, len(keys))
	for i, v := range keys {
		arr[i] = v
	}
	return arr
}

func init() {
	mp = make(map[string]cache.Cache)
	cache.Register(cache.DriverTypeRedis, GetRedisCache)
}
