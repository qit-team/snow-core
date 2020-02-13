package loggerfactory

import (
	"github.com/qit-team/snow-core/log/customlogger"
	"github.com/qit-team/snow-core/log/loggerinterface"
)

// save logger instance avoid create mutilple times
var loggerInstaces map[string]interface{}

func init() {
	loggerInstaces = map[string]interface{}{}
}

// GetLogger get logger implementation, the fileName is the file where log will write to and the first value of args is
// the logger interface implementation, if not found, return nil.
func GetLogger(fileName string, args ...string) loggerinterface.LoggerInterface {
	loggerName := ""
	if len(args) > 0 {
		loggerName = args[0]
	}
	if loggerName == "" {
		loggerName = "customlogger"
	}
	switch loggerName {
	case "customlogger":
		key := loggerName + "-" + fileName
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
