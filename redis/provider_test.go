package redis

import (
	"testing"
	"github.com/qit-team/snow-core/config"
)

func Test_getSingleton(t *testing.T) {
	c := getSingleton("", false)
	if c != nil {
		t.Error("client is not equal nil")
		return
	}
}

func TestProvider(t *testing.T) {
	err := Pr.Register("redis", config.RedisConfig{})
	if err == nil {
		t.Error(err)
		return
	}

	conf := config.RedisConfig{
		Master: config.RedisBaseConfig{
			Host: "127.0.0.1",
			Port: 6379,
		},
	}

	err = Pr.Register("redis", conf, true)
	if err != nil {
		t.Error(err)
		return
	}

	arr := Pr.Provides()
	if !(len(arr) == 1 && arr[0] == "redis") {
		t.Errorf("Provides is not match. %v", arr)
		return
	}

	err = Pr.Register("redis1", conf)
	if err != nil {
		t.Error(err)
		return
	}

	arr = Pr.Provides()
	if !(len(arr) == 2 && arr[1] == "redis1" || arr[1] == "redis") {
		t.Errorf("Provides is not match. %v", arr)
		return
	}

	c := GetRedis()
	if c == nil {
		t.Error("client is equal nil")
		return
	}

	c1 := GetRedis("redis1")
	if c1 == nil {
		t.Error("client is equal nil")
		return
	}

	defer func() {
		if e := recover(); e != "redis di_name:redis2 not exist" {
			t.Error("not panic")
		}
	}()
	GetRedis("redis2")

	err = Pr.Close()
	if err != nil {
		t.Error(err)
		return
	}
}
