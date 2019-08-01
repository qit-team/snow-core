package queue

import (
	"sync"
	"fmt"
)

const (
	DriverTypeRedis  = "redis"
	DriverTypeAliMns = "ali_mns"
)

var (
	drivers map[string]Instance
	mu      sync.RWMutex
)

type Instance func(diName string) Queue

func Register(driverType string, driver Instance) {
	if driver == nil {
		panic("queue.Register driver is nil")
	}
	mu.Lock()
	defer mu.Unlock()

	if _, ok := drivers[driverType]; ok {
		panic("queue.Register called twice for driver " + driverType)
	}
	drivers[driverType] = driver
}

//获取Queue对象
func GetQueue(diName string, driverType string) (q Queue) {
	mu.RLock()
	instanceFunc, ok := drivers[driverType]
	mu.RUnlock()
	if !ok {
		panic(fmt.Sprintf("queue.GetQueue unknown driver %s", driverType))
	}
	q = instanceFunc(diName)
	if q == nil {
		panic(fmt.Sprintf("queue.GetQueue unknown diName %s", diName))
	}
	return
}

func init() {
	drivers = make(map[string]Instance)
}
