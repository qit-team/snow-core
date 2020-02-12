package logger

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/qit-team/snow-core/config"
	"github.com/qit-team/snow-core/http/ctxkit"
)

var contextTest, contextTest1 *gin.Context

func init() {
	contextTest = &gin.Context{}
	contextTest1 = &gin.Context{}
}

func Test_getSingleton(t *testing.T) {
	c := getSingleton("", false)
	if c != nil {
		t.Error("client is not equal nil")
		return
	}
}

func TestProvider(t *testing.T) {
	err := Pr.Register("logger", config.LogConfig{})
	if err == nil {
		t.Error(err)
		return
	}

	conf := config.LogConfig{
		Handler: "file",
		Level:   "info",
		Dir:     "../../",
	}

	err = Pr.Register("logger", conf, true)
	if err != nil {
		t.Error(err)
		return
	}

	// test generate trace id
	traceId, err := ctxkit.GenerateTraceId(contextTest)

	if err != nil {
		t.Error("generateTraceIdError", err, traceId)
	}

	// 对context设置traceId
	ctxkit.SetTraceId(contextTest, traceId)
	temp := ctxkit.GetTraceId(contextTest)
	fmt.Println("=======test_temp:", temp)
	Info(contextTest, "========testTraceId:levelInfo=====")
	Error(contextTest, "========testTraceId:levelError=====")
	Warn(contextTest, "========testTraceId:levelWarn=====")
	Debug(contextTest, "========testTraceId:levelDebug=====")
	Trace(contextTest, "========testTraceId:levelTrace=====")
	//Fatal(contextTest, "========testTraceId:levelFatal=====")

	Info(nil, "================")

	// 新的context，确保第一次记录log，会在context中种下traceId
	Info(contextTest1, "========testTraceId111:levelInfo=====")
	Error(contextTest1, "========testTraceId111:levelError=====")
	Warn(contextTest1, "========testTraceId111:levelWarn=====")
	Debug(contextTest1, "========testTraceId111:levelDebug=====")
	// 调用panic会导致go test fail
	//Panic(contextTest1, "========testTraceId111:levelPanic=====")

	arr := Pr.Provides()
	if !(len(arr) == 1 && arr[0] == "logger") {
		t.Errorf("Provides is not match. %v", arr)
		return
	}

	err = Pr.Register("logger1", conf)
	if err != nil {
		t.Error(err)
		return
	}

	arr = Pr.Provides()
	if !(len(arr) == 2 && arr[1] == "logger1" || arr[1] == "logger") {
		t.Errorf("Provides is not match. %v", arr)
		return
	}

	c := GetLogger()
	if c == nil {
		t.Error("client is equal nil")
		return
	}

	c1 := GetLogger("logger1")
	if c1 == nil {
		t.Error("client is equal nil")
		return
	}

	defer func() {
		if e := recover(); e != "logger di_name:logger2 not exist" {
			t.Error("not panic")
		}
	}()
	GetLogger("logger2")

	err = Pr.Close()
	if err != nil {
		t.Error(err)
		return
	}

}

func TestNewWithField(t *testing.T) {
	// 测试NewWithField && BatchNewWithField方法
	conf := config.LogConfig{
		Handler: "file",
		Level:   "info",
		Dir:     "../../",
	}

	defer func() {
		if e := recover(); e != nil {
			t.Error("test NewWithField panic")
		}
	}()
	err := Pr.Register("logger", conf, true)
	if err != nil {
		t.Error(err)
		return
	}

	Info(nil, "===TestNewWithField", NewWithField("data", "snow"))

	logInfo := map[string]interface{}{
		"url":    "testUrl",
		"params": "snow",
		"num":    100,
	}

	// Info(nil, "===TestBatchNewWithField", BatchNewWithField(logInfo), "asdfasdfasdasdfasd")

	Info(nil, "===TestBatchNewWithField", logInfo, "asdfasdfasdasdfasd")

}

func TestWithFileName(t *testing.T) {
	// 测试NewWithField && BatchNewWithField方法
	conf := config.LogConfig{
		Handler:  "file",
		Level:    "info",
		Dir:      "../../",
		FileName: "TestWithFileName",
	}

	defer func() {
		if e := recover(); e != nil {
			t.Error("test NewWithField panic")
		}
	}()

	err := Pr.Register("logger", conf, true)
	if err != nil {
		t.Error(err)
		return
	}

	Info(nil, "asdfasdfasdf", map[string]interface{}{
		"c": 123,
	})

	// GetLoggerWithFileName("nihao").WithContext(context.Background()).Info("asdfaosdfaosdihfaposd")
	// GetLoggerWithFileName("buhao").WithField("a", "b").Info(map[string]interface{}{
	// 	"c": 123,
	// })
}
