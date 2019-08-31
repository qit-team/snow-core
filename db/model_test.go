package db

import (
	"fmt"
	"github.com/qit-team/snow-core/config"
	"testing"

	//go test时需要开启
	_ "github.com/go-sql-driver/mysql"
)

type bannerModel struct {
	Model
}

func init() {
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

	err := Pr.Register("db", dbConf, true)
	if err != nil {
		fmt.Println(err)
	}

	engineGroup = GetDb()
}

func TestGetOne(t *testing.T) {
	model := new(bannerModel)
	ret := new(Banner)
	id := 1
	_, err := model.GetOne(id, ret)
	if err != nil {
		t.Errorf("getOne error: %v", err)
		return
	}
	fmt.Println("getOne.Ret", ret)
}

func TestGetMulti(t *testing.T) {
	model := new(bannerModel)
	ret := make([]*Banner, 0)
	var idList = []interface{}{1, 2}
	err := model.GetMulti(idList, &ret)
	if err != nil {
		t.Errorf("getMulti error: %v", err)
		return
	}
	for _, v := range ret {
		fmt.Println("getMulti.ItemRet", v)
	}

	// 验证异常处理if分支
	var idListErr []interface{}
	err = model.GetMulti(idListErr, &ret)
	fmt.Println("getMulti.CheckExceptionBranch:", err)
}

func TestInsert(t *testing.T) {

	model := new(bannerModel)
	banner := new(Banner)
	banner.Id = 4
	banner.ImageUrl = "img666"
	banner.Pid = 66666
	banner.Title = "test insert"

	_, err := model.Insert(banner)
	if err != nil {
		t.Errorf("Insert error: %v", err)
		return
	}
	fmt.Println("Insert.Id", banner.Id)

	// 插入数据 为了测试批量删除功能，见函数：TestDeleteMulti
	banner.Id = 5
	model.Insert(banner)

	banner.Id = 6
	model.Insert(banner)
}

func TestUpdate(t *testing.T) {
	model := new(bannerModel)
	banner := new(Banner)
	banner.ImageUrl = ""
	banner.Pid = 77777
	banner.Title = "test update"
	var id = 7
	// 注意：直接用默认的update对上面的ImageUrl字段不会更新
	_, err := model.Update(id, banner)
	if err != nil {
		t.Errorf("Update error: %v", err)
		return
	}
	fmt.Println("Update.success")

	id = 8
	banner.Pid = 888
	banner.ImageUrl = ""
	banner.Title = ""
	// xorm默认对更新字段数据为""的不会执行，需要加mustColumns，这样保证为空的数据字段也能更新，详情搜索xorm手册
	_, err = model.Update(id, banner, "img_url", "title")

	if err != nil {
		t.Errorf("Update mustColumns error: %v", err)
		return
	}
	fmt.Println("Update mustColumns.success")
}

func TestDelete(t *testing.T) {
	model := new(bannerModel)
	banner := new(Banner)
	id := 4
	ret, err := model.Delete(id, banner)

	if err != nil {
		t.Errorf("Delete error: %v", err)
		return
	}
	fmt.Println("Delete.ret", ret)
}

func TestDeleteMulti(t *testing.T) {
	model := new(bannerModel)
	banner := new(Banner)
	var id = []interface{}{5, 6}
	// 批量删除id 为5，6的  数据来源参考TestInsert
	ret, err := model.DeleteMulti(id, banner)

	if err != nil {
		t.Errorf("DeleteMulti error: %v", err)
		return
	}
	fmt.Println("DeleteMulti.ret", ret)

	// 测试参数为空的异常分支
	var idErr []interface{}
	_, err = model.DeleteMulti(idErr, banner)
	fmt.Println("DeleteMulti.CheckExceptionBranch.ret", err)
}

func TestGetList(t *testing.T) {
	model := new(bannerModel)
	banner := make([]*Banner, 0)

	sql := "status > ? and status < ? and pid = ?"
	var values = []interface{}{"1", "5", 10010}
	err := model.GetList(&banner, sql, values)
	if err != nil {
		t.Errorf("Getlist error: %v", err)
		return
	}
	for _, v := range banner {
		fmt.Println("GetList.ret", v)
	}

	// 测试其他if分支 覆盖getList所有代码
	banner1 := make([]*Banner, 0)

	sql = "status >= ? and status <= ?"
	var valuesTest = []interface{}{"1", "7"}
	err = model.GetList(&banner1, sql, valuesTest, []int{3, 3}, "pid desc")
	if err != nil {
		t.Errorf("GetlistLimitAndOrderBranch error: %v", err)
		return
	}
	for _, v := range banner1 {
		fmt.Println("GetlistLimitAndOrderBranch.ret", v)
	}
}

func TestProvider_Close(t *testing.T) {
	// 关闭链接，此时再执行sql都无法执行会报 sql: database is closed， 所以在sql执行完之后做close操作
	err := Pr.Close()

	if err != nil {
		t.Error("Close Fail")
		return
	}
}
