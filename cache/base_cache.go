package cache

import (
	"github.com/qit-team/snow-core/utils"
	"context"
	"github.com/qit-team/snow-core/redis"
)

const (
	DefaultDiName     = redis.SingletonMain
	DefaultDriverType = DriverTypeRedis
	DefaultPrefix     = ""    //默认缓存key前缀
	DefaultTTL        = 86400 //默认缓存时间
)

//缓存基类
type BaseCache struct {
	cache      Cache
	DiName     string //缓存依赖的实例别名
	Prefix     string //缓存key前缀
	DriverType string //缓存驱动
	ttl        int    //缓存时间
	ttlIsSet   bool   //避免TTL被设置过为0时，仍使用默认值的情况
}

//补全key
func (m *BaseCache) key(key string) string {
	return m.Prefix + key
}

//批量补全
func (m *BaseCache) keys(keys ...string) []string {
	arr := make([]string, len(keys))
	for i, key := range keys {
		arr[i] = m.key(key)
	}
	return arr
}

//去除前缀
func (m *BaseCache) removePrefix(key string) string {
	l := len(m.Prefix)
	return utils.Substr(key, l, len(key)-l)
}

func (m *BaseCache) GetPrefixOrDefault() string {
	if m.Prefix != "" {
		return m.Prefix
	} else {
		return DefaultPrefix
	}
}

func (m *BaseCache) GetDiNameOrDefault() string {
	if m.DiName != "" {
		return m.DiName
	} else {
		return DefaultDiName
	}
}

func (m *BaseCache) GetDriverTypeOrDefault() string {
	if m.DriverType != "" {
		return m.DriverType
	} else {
		return DefaultDriverType
	}
}

func (m *BaseCache) SetTTL(ttl int) {
	m.ttlIsSet = true
	m.ttl = ttl
}

func (m *BaseCache) GetTTLOrDefault() int {
	if m.ttlIsSet {
		return m.ttl
	} else {
		return DefaultTTL
	}
}

func (m *BaseCache) getTTL(ttl ...int) int {
	if len(ttl) > 0 {
		return ttl[0]
	} else {
		return m.GetTTLOrDefault()
	}
}

func (m *BaseCache) Get(ctx context.Context, key string) (interface{}, error) {
	key = m.key(key)
	return m.GetCache().Get(ctx, key)
}

func (m *BaseCache) Set(ctx context.Context, key string, value interface{}, ttl ...int) (bool, error) {
	key = m.key(key)
	return m.GetCache().Set(ctx, key, value, m.getTTL(ttl...))
}

func (m *BaseCache) GetMulti(ctx context.Context, keys ...string) (map[string]interface{}, error) {
	keys = m.keys(keys...)
	items, err := m.GetCache().GetMulti(ctx, keys...)
	if err != nil {
		return nil, err
	}

	m2 := make(map[string]interface{})
	for key, val := range items {
		m2[m.removePrefix(key)] = val
	}
	return m2, nil
}

func (m *BaseCache) SetMulti(ctx context.Context, items map[string]interface{}, ttl ...int) (bool, error) {
	arr := make(map[string]interface{})
	for key, value := range items {
		key = m.key(key)
		arr[key] = value
	}
	return m.GetCache().SetMulti(ctx, arr, m.getTTL(ttl...))
}

func (m *BaseCache) Delete(ctx context.Context, key string) (bool, error) {
	key = m.key(key)
	return m.GetCache().Delete(ctx, key)
}

func (m *BaseCache) DeleteMulti(ctx context.Context, keys ...string) (bool, error) {
	keys = m.keys(keys...)
	return m.GetCache().DeleteMulti(ctx, keys...)
}

func (m *BaseCache) Expire(ctx context.Context, key string, ttl ...int) (bool, error) {
	key = m.key(key)
	return m.GetCache().Expire(ctx, key, m.getTTL(ttl...))
}

func (m *BaseCache) IsExist(ctx context.Context, key string) (bool, error) {
	key = m.key(key)
	return m.GetCache().IsExist(ctx, key)
}

//获取缓存类
func (m *BaseCache) GetCache() Cache {
	//不使用once.Done是因为会有多种cache实例
	diName := m.GetDiNameOrDefault()
	driverType := m.GetDriverTypeOrDefault()
	return GetCache(diName, driverType)
}
