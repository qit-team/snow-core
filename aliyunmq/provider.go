package aliyunmq

import (
	"errors"
	"fmt"
	"github.com/aliyunmq/mq-http-go-sdk"
	"github.com/qit-team/snow-core/config"
	"github.com/qit-team/snow-core/helper"
	"github.com/qit-team/snow-core/kernel/container"
	"sync"
)

const (
	SingletonMain = "aliyun_mq"
)

var Pr *provider

func init() {
	Pr = new(provider)
	Pr.mp = make(map[string]interface{})
}

type provider struct {
	mu sync.RWMutex
	mp map[string]interface{} //配置
	dn string                 //default name
}

/**
 * @param string 依赖注入别名 必选
 * @param config.LogConfig 配置 必选
 * @param bool 是否启用懒加载 可选
 */
func (p *provider) Register(args ...interface{}) (err error) {
	diName, lazy, err := helper.TransformArgs(args...)
	if err != nil {
		return
	}

	conf, ok := args[1].(config.AliyunMqConfig)
	if !ok {
		return errors.New("args[1] is not config.MnsConfig")
	}

	p.mu.Lock()
	p.mp[diName] = args[1]
	if len(p.mp) == 1 {
		p.dn = diName
	}
	p.mu.Unlock()

	if !lazy {
		_, err = setSingleton(diName, conf)
	}
	return
}

//注册过的别名
func (p *provider) Provides() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return helper.MapToArray(p.mp)
}

//释放资源
func (p *provider) Close() error {
	return nil
}

//注入单例
func setSingleton(diName string, conf config.AliyunMqConfig) (ins mq_http_sdk.MQClient, err error) {
	ins, err = NewAliyunMqClient(conf)
	if err == nil {
		container.App.SetSingleton(diName, ins)
	}
	return
}

//获取单例
func getSingleton(diName string, lazy bool) mq_http_sdk.MQClient {
	rc := container.App.GetSingleton(diName)
	if rc != nil {
		return rc.(mq_http_sdk.MQClient)
	}
	if lazy == false {
		return nil
	}

	Pr.mu.RLock()
	conf, ok := Pr.mp[diName].(config.AliyunMqConfig)
	Pr.mu.RUnlock()
	if !ok {
		panic(fmt.Sprintf("aliyun_mq di_name:%s not exist", diName))
	}

	ins, err := setSingleton(diName, conf)
	if err != nil {
		panic(fmt.Sprintf("aliyun_mq di_name:%s err:%s", diName, err.Error()))
	}
	return ins
}

//外部通过注入别名获取资源，解耦资源的关系
func GetAliyunMq(args ...string) mq_http_sdk.MQClient {
	diName := helper.GetDiName(Pr.dn, args...)
	return getSingleton(diName, true)
}
