package command

import (
	"sync"
	"errors"
)

var (
	ErrUnknownName = errors.New("unknown name")
)

//一次性任务脚本
type Command struct {
	mu        sync.RWMutex
	container map[string]func()
}

//new实例
func New() *Command {
	c := new(Command)
	c.container = make(map[string]func())
	return c
}

//绑定name与函数的关系
func (c *Command) AddFunc(name string, f func()) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.container[name] = f
}

//通过name执行函数
func (c *Command) Execute(name string) (err error) {
	c.mu.RLock()
	f, ok := c.container[name]
	c.mu.RUnlock()
	if ok {
		f()
	} else {
		panic(ErrUnknownName.Error())
	}
	return
}
