package aliyunmq

import (
	"fmt"
	"github.com/qit-team/snow-core/config"
	"testing"
)

func Test_getSingleton(t *testing.T) {
	c := getSingleton("", false)
	if c != nil {
		t.Error("client is not equal nil")
		return
	}
}

func TestProvider(t *testing.T) {
	err := Pr.Register("aliyun_mq", config.AliyunMqConfig{}, true)
	if err != nil {
		t.Error(err)
		return
	}

	arr := Pr.Provides()
	if !(len(arr) == 1 && arr[0] == "aliyun_mq") {
		t.Errorf("Provides is not match. %v", arr)
		return
	}

	err = Pr.Register("aliyun_mq1", config.AliyunMqConfig{})
	if err != nil {
		t.Error(err)
		return
	}

	arr = Pr.Provides()
	if !(len(arr) == 2 && arr[1] == "aliyun_mq" || arr[1] == "aliyun_mq1") {
		t.Errorf("Provides is not match. %v", arr)
		return
	}

	err = Pr.Close()
	if err != nil {
		t.Error(err)
		return
	}

	c := GetAliyunMq()
	fmt.Println("providers.GetAliyunMq:", c)

	if c != nil {
		t.Error("client is not equal nil")
		return
	}
}
