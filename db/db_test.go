package db

import (
	"testing"
	"github.com/qit-team/snow-core/config"
	"github.com/go-xorm/xorm"
	"fmt"
	//go test时需要开启
	_ "github.com/go-sql-driver/mysql"
)

var engineGroup *xorm.EngineGroup

/**
 * Banner实体
 */
type Banner struct {
	Id       int64  `xorm:"pk autoincr"`
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
	m := config.DbBaseConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "root",
		Password: "123456",
		DBName:   "test",
	}
	dbConf := config.DbConfig{
		Driver: "mysql",
		Master: m,
	}

	err := Pr.Register("db", dbConf, true)
	if err != nil {
		fmt.Println(err)
	}

	engineGroup = GetDb()
}

func TestGet(t *testing.T) {
	banner := new(Banner)
	engineGroup.ShowSQL(true)
	_, err := engineGroup.ID(1).Get(banner)

	if err != nil {
		t.Errorf("get error: %v", err)
		return
	}

	fmt.Println(banner)
}
