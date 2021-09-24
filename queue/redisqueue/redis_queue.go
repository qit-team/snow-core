package redisqueue

import (
	"context"
	"errors"
	goredis "github.com/go-redis/redis/v8"
	"github.com/qit-team/snow-core/queue"
	"github.com/qit-team/snow-core/redis"
	"sync"
)

var (
	mp map[string]queue.Queue
	mu sync.RWMutex
)

type RedisQueue struct {
	client *goredis.Client
}

//new实例
func newRedisQueue(diName string) queue.Queue {
	m := new(RedisQueue)
	m.client = redis.GetRedis(diName)
	return m
}

//单例模式
func GetRedisQueue(diName string) queue.Queue {
	key := diName
	mu.RLock()
	q, ok := mp[key]
	mu.RUnlock()
	if ok {
		return q
	}

	q = newRedisQueue(diName)
	mu.Lock()
	mp[key] = q
	mu.Unlock()
	return q
}

/**
 * 队列消息入队
 */
func (m *RedisQueue) Enqueue(ctx context.Context, key string, message string, args ...interface{}) (bool, error) {
	//redis暂时不要延迟和优先级
	_, err := m.client.RPush(ctx, key, message).Result()
	if err != nil {
		return false, err
	}
	return true, err
}

/**
 * 队列消息出队
 */
func (m *RedisQueue) Dequeue(ctx context.Context, key string, args ...interface{}) (message string, tag string, token string, dequeueCount int64, err error) {
	// redis 出队次数暂用1 目前不支持统计这个次数
	dequeueCount = 0
	message, err = m.client.LPop(ctx, key).Result()
	if err == goredis.Nil {
		err = nil
		message = ""
	}
	return
}

/**
 * 确认消息接收 redis暂时用不到
 */
func (m *RedisQueue) AckMsg(ctx context.Context, key string, token string, args ...interface{}) (bool, error) {
	return true, nil
}

/**
 * 队列消息入队
 */
func (m *RedisQueue) BatchEnqueue(ctx context.Context, key string, messages []string, args ...interface{}) (bool, error) {
	//redis暂时不要延迟和优先级
	if len(messages) == 0 {
		return false, errors.New("messages is empty")
	}
	_, err := m.client.RPush(ctx, key, arrayStringToInterface(messages)...).Result()
	if err != nil {
		return false, err
	}
	return true, err
}

func arrayStringToInterface(arr []string) []interface{} {
	newArr := make([]interface{}, len(arr))
	for k, v := range arr {
		newArr[k] = v
	}
	return newArr
}

func init() {
	mp = make(map[string]queue.Queue)
	queue.Register(queue.DriverTypeRedis, GetRedisQueue)
}
