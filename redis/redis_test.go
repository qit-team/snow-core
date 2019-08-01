package redis

import (
	"testing"
	"reflect"
	"strconv"
	"github.com/qit-team/snow-core/config"
	"time"
)

var conf config.RedisConfig

func init() {
	conf = config.RedisConfig{
		Master: config.RedisBaseConfig{
			Host: "127.0.0.1",
			Port: 6379,
		},
	}
}

func TestGetSet(t *testing.T) {
	_, err := NewRedisClient(config.RedisConfig{})
	if err == nil {
		t.Error("redis config donot check")
		return
	}

	client, err := NewRedisClient(conf)
	if err != nil {
		t.Error("client init failed")
		return
	}

	value := 11
	res, _ := client.Set("hts", value)
	t.Log(res, reflect.TypeOf(res))
	if res == false {
		t.Error("set error")
		return
	}

	res1, _ := client.Get("hts")
	t.Log(res1, reflect.TypeOf(res1))
	if res1 == "" {
		t.Error("get error")
		return
	} else if res1 != strconv.Itoa(value) {
		t.Error("not same")
		return
	}
}

func Test_genRedisConfig(t *testing.T) {
	conf := config.RedisBaseConfig{
		Host: "127.0.0.1",
		Port: 6379,
	}
	newConf := genRedisConfig(conf)
	if newConf.Host != conf.Host || newConf.Port != conf.Port || newConf.DB != conf.DB {
		t.Error("genRedisConfig failed")
		return
	}
}

func Test_genOptions(t *testing.T) {
	conf := config.RedisOptionConfig{
		MaxConns:    64,
		Wait:        true,
		IdleTimeout: 3,
	}
	newConf := genOptions(conf)
	if newConf.MaxIdle != 0 || newConf.Wait != conf.Wait || newConf.MaxActive != conf.MaxConns {
		t.Error("genOptions failed")
		return
	} else if newConf.IdleTimeout != 3*time.Second || newConf.ConnectTimeout != 0*time.Second {
		t.Error("genOptions failed")
		return
	}
}
