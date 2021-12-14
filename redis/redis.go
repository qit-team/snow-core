package redis

import (
	"errors"
	"fmt"
	goredis "github.com/go-redis/redis/v8"
	"github.com/qit-team/snow-core/config"
	"time"
)

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}
type Options struct {
	MaxIdle        int
	MaxActive      int
	Wait           bool
	IdleTimeout    time.Duration
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

//redis连接池实例，不对外暴露，通过redis_service_provider实现依赖注入和资源获取
func NewRedisClient(redisConf config.RedisConfig) (*goredis.Client, error) {
	if redisConf.Master.Host == "" {
		return nil, errors.New("redis config is empty")
	}

	rdb := goredis.NewClient(&goredis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConf.Master.Host, redisConf.Master.Port),
		Password: redisConf.Master.Password,
		DB:       redisConf.Master.DB,
	})
	return rdb, nil
}

//redis连接池实例，不对外暴露，通过redis_service_provider实现依赖注入和资源获取
func NewClusterRedisClient(redisConf config.RedisConfig) (*goredis.ClusterClient, error) {
	if redisConf.Master.Host == "" {
		return nil, errors.New("redis config is empty")
	}

	addrs := []string{}
	addrs = append(addrs, fmt.Sprintf("%s:%d", redisConf.Master.Host, redisConf.Master.Port))
	for _, slave := range redisConf.Slaves {
		addrs = append(addrs, fmt.Sprintf("%s:%d", slave.Host, slave.Port))
	}
	rdb := goredis.NewClusterClient(&goredis.ClusterOptions{
		Addrs:    addrs,
		Password: redisConf.Master.Password,
	})
	return rdb, nil
}

func genRedisConfig(c config.RedisBaseConfig) RedisConfig {
	return RedisConfig{
		Host:     c.Host,
		Port:     c.Port,
		Password: c.Password,
		DB:       c.DB,
	}
}

func genOptions(c config.RedisOptionConfig) Options {
	return Options{
		MaxIdle:        c.MaxIdle,
		MaxActive:      c.MaxConns,
		Wait:           c.Wait,
		IdleTimeout:    c.IdleTimeout * time.Second,
		ConnectTimeout: c.ConnectTimeout * time.Second,
		ReadTimeout:    c.ReadTimeout * time.Second,
		WriteTimeout:   c.WriteTimeout * time.Second,
	}
}
