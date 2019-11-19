package aliyunmq

import (
	"fmt"
	"github.com/qit-team/snow-core/config"
	"io/ioutil"
	"strings"
	"testing"
)

func TestNewAliyunMqClient(t *testing.T) {

	conf := getConfig()
	fmt.Println("aliyun_mq config:", conf)
	c, err := NewAliyunMqClient(conf)
	if err != nil {
		t.Error(err)
		return
	} else if c == nil {
		t.Error("client is nil")
		return
	}
}

func TestNewAliyunMqClient2(t *testing.T) {
	conf := config.AliyunMqConfig{
		EndPoint:  "",
		AccessKey: "",
		SecretKey: "",
	}
	_, err := NewAliyunMqClient(conf)
	if err == nil {
		t.Error("invalid config must return err")
	}
}

func getConfig() config.AliyunMqConfig {
	//需要自己在文件填好配置
	bs, err := ioutil.ReadFile("../.env.aliyunmq")

	conf := config.AliyunMqConfig{}
	if err == nil {
		str := string(bs)
		arr := strings.Split(str, "\n")
		if len(arr) >= 3 {
			conf.EndPoint = arr[0]
			conf.AccessKey = arr[1]
			conf.SecretKey = arr[2]
		}
	}
	return conf
}
