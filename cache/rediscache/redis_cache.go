package rediscache

import (
	"context"
	goredis "github.com/go-redis/redis/v8"
	"github.com/qit-team/snow-core/cache"
	"github.com/qit-team/snow-core/redis"
	"sync"
	"time"
)

var (
	mp map[string]cache.Cache
	mu sync.RWMutex
)

type RedisCache struct {
	client *goredis.Client
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
	value, err :=  c.client.Get(ctx,key).Result()
	if err == goredis.Nil {
		return "", nil
	}
	return value, err
}

func (c *RedisCache) GetMulti(ctx context.Context, keys ...string) (map[string]interface{}, error) {
	values, err := c.client.MGet(ctx,keys...).Result()
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
	_,err := c.client.SetEX(ctx,key, value, time.Duration(t)*time.Second).Result()
	if err != nil {
		return false,nil
	}
	return true,nil
}

func (c *RedisCache) SetMulti(ctx context.Context, items map[string]interface{}, ttl ...int) (bool, error) {
	arr := make([]interface{}, 0)
	for key, value := range items {
		arr = append(arr, key, value)
	}
	_, err := c.client.MSet(ctx,arr...).Result()
	if err != nil {
		return false, err
	}

	t := cache.GetTTLOrDefault(ttl...)
	if t > 0 {
		t64 := int64(t)
		for key, _ := range items {
			c.client.Expire(ctx,key, time.Duration(t64)*time.Second)
		}
	}
	return true, nil
}

func (c *RedisCache) Delete(ctx context.Context, key string) (bool, error) {
	res, err := c.client.Del(ctx,key).Result()
	return res > 0, err
}

func (c *RedisCache) DeleteMulti(ctx context.Context, keys ...string) (bool, error) {
	res, err := c.client.Del(ctx,keys...).Result()
	return res > 0, err
}

func (c *RedisCache) Expire(ctx context.Context, key string, ttl ...int) (bool, error) {
	t := cache.GetTTLOrDefault(ttl...)
	return c.client.Expire(ctx,key, time.Duration(t)*time.Second).Result()
}

func (c *RedisCache) IsExist(ctx context.Context, key string) (bool, error) {
	num, err := c.client.Exists(ctx,key).Result()
	return num == 1, err
}

func convert(keys []string) []interface{} {
	arr := make([]interface{}, len(keys))
	for i, v := range keys {
		arr[i] = v
	}
	return arr
}

func (c *RedisCache) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	newVal, err := c.client.IncrBy(ctx,key, value).Result()
	if err != nil {
		return 0, err
	}
	return newVal, err
}

func (c *RedisCache) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	newVal, err := c.client.DecrBy(ctx,key, value).Result()
	if err != nil {
		return 0, err
	}
	return newVal, err
}

func init() {
	mp = make(map[string]cache.Cache)
	cache.Register(cache.DriverTypeRedis, GetRedisCache)
}
