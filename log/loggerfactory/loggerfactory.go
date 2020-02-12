package loggerfactory

import (
	"github.com/qit-team/snow-core/log/customlogger"
	"github.com/qit-team/snow-core/log/loggerinterface"
)

var loggerInstace map[string]interface{}

func init() {
	loggerInstace = map[string]interface{}{}
}

func GetLogger(fileName string, loggerType string) loggerinterface.LoggerInterface {
	if loggerType == "" {
		loggerType = "customlogger"
	}
	switch loggerType {
	case "customlogger":
		key := loggerType + "-" + fileName
		if value, ok := loggerInstace[key]; ok {
			return value.(*customlogger.CustomLogger)
		}
		instance := new(customlogger.CustomLogger)
		instance.FileName = fileName
		loggerInstace[key] = instance
		return instance
	default:
		return nil
	}
}
