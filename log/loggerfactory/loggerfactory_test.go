package loggerfactory

import (
	"testing"

	"github.com/qit-team/snow-core/config"
	"github.com/qit-team/snow-core/log/logger"
)

func TestLoggerFactory(t *testing.T) {
	conf := config.LogConfig{
		Handler:  "file",
		Level:    "info",
		Dir:      "../../",
		FileName: "TestLoggerFactory",
	}

	defer func() {
		if e := recover(); e != nil {
			t.Error("test TestLoggerFactory panic")
		}
	}()

	err := logger.Pr.Register("logger", conf, true)
	if err != nil {
		t.Error(err)
		return
	}

	fileName := "loggerfactory"
	GetLogger(fileName).Info(nil, "Test logger factory", map[string]interface{}{
		"name": "liou",
	})
}
