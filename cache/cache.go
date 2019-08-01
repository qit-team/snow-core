package cache

import (
	"sync"
	"fmt"
)

const (
	DriverTypeRedis = "redis"
)

var (
	drivers map[string]Instance
	mu      sync.RWMutex
)

type Instance func(diName string) Cache

func Register(driverType string, driver Instance) {
	if driver == nil {
		panic("cache.Register driver is nil")
	}
	mu.Lock()
	defer mu.Unlock()

	if _, ok := drivers[driverType]; ok {
		panic("cache.Register called twice for driver " + driverType)
	}
	drivers[driverType] = driver
}

// args columns: TTL int
func GetCache(diName string, driverType string) (q Cache) {
	mu.RLock()
	instanceFunc, ok := drivers[driverType]
	mu.RUnlock()
	if !ok {
		panic(fmt.Sprintf("cache.GetCache unknown driver %s", driverType))
	}
	q = instanceFunc(diName)
	if q == nil {
		panic(fmt.Sprintf("cache.GetCache unknown diName %s", diName))
	}
	return
}

//获取TTL时间
func GetTTLOrDefault(ttl ...int) (t int) {
	if len(ttl) > 0 {
		t = ttl[0]
	} else {
		t = DefaultTTL
	}
	return
}

func init() {
	drivers = make(map[string]Instance)
}
