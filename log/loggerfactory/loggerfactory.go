package loggerfactory

import (
	"github.com/qit-team/snow-core/log/customlogger"
	"github.com/qit-team/snow-core/log/loggerinterface"
)

var loggerInstaces map[string]interface{}

func init() {
	loggerInstaces = map[string]interface{}{}
}

func GetLogger(fileName string, loggerType string) loggerinterface.LoggerInterface {
	if loggerType == "" {
		loggerType = "customlogger"
	}
	switch loggerType {
	case "customlogger":
		key := loggerType + "-" + fileName
		if value, ok := loggerInstaces[key]; ok {
			return value.(*customlogger.CustomLogger)
		}
		instance := new(customlogger.CustomLogger)
		instance.FileName = fileName
		loggerInstaces[key] = instance
		return instance
	default:
		return nil
	}
}
