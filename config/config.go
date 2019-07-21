package config

import "time"

type RedisBaseConfig struct {
	Host     string
	Port     int
	Password string
	DB       int //第几个库，默认0
}

type RedisOptionConfig struct {
	MaxIdle        int
	MaxConns       int
	Wait           bool
	IdleTimeout    time.Duration
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

type RedisConfig struct {
	Master RedisBaseConfig
	Slaves []RedisBaseConfig
	Option RedisOptionConfig
}

type DbBaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type DbOptionConfig struct {
	MaxIdle        int
	MaxConns       int
	IdleTimeout    time.Duration
	ConnectTimeout time.Duration
	Charset        string
}

type DbConfig struct {
	Driver string //驱动类型，目前支持mysql、postgres、mssql、sqlite3
	Master DbBaseConfig
	Slaves []DbBaseConfig
	Option DbOptionConfig
}

type MnsConfig struct {
	Url             string
	AccessKeyId     string
	AccessKeySecret string
}

type LogConfig struct {
	Handler  string
	Level    string
	Dir      string
	FileName string
}

type ApiConfig struct {
	Host string
	Port int
}
