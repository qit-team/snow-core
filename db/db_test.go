package db

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/qit-team/snow-core/config"
	"testing"
	//go test时需要开启
	_ "github.com/go-sql-driver/mysql"
)

var engineGroup *xorm.EngineGroup

/**
 * Banner实体
 */
type Banner struct {
	Id       int64 `xorm:"pk autoincr"`
	Pid      int
	Title    string
	ImageUrl string `xorm:"'img_url'"`
}

/**
 * 表名规则
 */
func (m *Banner) TableName() string {
	return "banner"
}

func init() {
	dbInit(true)
}

func dbInit(lazyBool bool) {
	m := config.DbBaseConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "root",
		Password: "Snow_123",
		DBName:   "test",
	}
	dbConf := config.DbConfig{
		Driver: "mysql",
		Master: m,
	}

	err := Pr.Register("db", dbConf, lazyBool)
	if err != nil {
		fmt.Println(err)
	}

	engineGroup = GetDb()
}

func TestGet(t *testing.T) {
	banner := new(Banner)
	// sql是否打印开关
	//engineGroup.ShowSQL(true)
	_, err := engineGroup.ID(1).Get(banner)

	if err != nil {
		t.Errorf("get error: %v", err)
		return
	}

	fmt.Println(banner)
}

func TestProvider_Provides(t *testing.T) {
	retList := Pr.Provides()
	if len(retList) == 0 {
		t.Error("Provides empty")
		return
	}

	for k, v := range retList {
		fmt.Println("Provides list", k, v)
	}
}
