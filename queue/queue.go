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
		panic("queue.Register called twice for driver type:" + driverType)
	}
	drivers[driverType] = driver
}

//获取Queue对象
func GetQueue(diName string, driverType string) (q Queue) {
	mu.RLock()
	instanceFunc, ok := drivers[driverType]
	if !ok {
		panic(fmt.Sprintf("queue.GetQueue driver type (%s) is not exist", driverType))
	}
	q = instanceFunc(diName)
	if q == nil {
		panic(fmt.Sprintf("queue.GetQueue driver (%s:%s) is nil", diName, driverType))
	}
	return
}

func init() {
	drivers = make(map[string]Instance)
}
