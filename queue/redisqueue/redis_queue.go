package redisqueue

import (
	"context"
	redis_pool "github.com/hetiansu5/go-redis-pool"
	"github.com/qit-team/snow-core/redis"
	"errors"
	"sync"
	"github.com/qit-team/snow-core/queue"
)

var (
	mp map[string]queue.Queue
	mu sync.RWMutex
)

type RedisQueue struct {
	client *redis_pool.ReplicaPool
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
	_, err := m.client.RPush(key, message)
	if err != nil {
		return false, err
	}
	return true, err
}

/**
 * 队列消息出队
 */
func (m *RedisQueue) Dequeue(ctx context.Context, key string) (message string, token string, err error) {
	message, err = m.client.LPop(key)
	if err == redis_pool.ErrNil {
		err = nil
		message = ""
	}
	return
}

/**
 * 确认消息接收 redis暂时用不到
 */
func (m *RedisQueue) AckMsg(ctx context.Context, key string, token string) (bool, error) {
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
	_, err := m.client.RPush(key, arrayStringToInterface(messages)...)
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
