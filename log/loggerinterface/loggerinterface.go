package loggerinterface

import "context"

// LoggerInterface 日志接口
type LoggerInterface interface {
	Trace(c context.Context, msg ...interface{})

	Debug(c context.Context, msg ...interface{})

	Info(c context.Context, msg ...interface{})

	Warn(c context.Context, msg ...interface{})

	Error(c context.Context, msg ...interface{})

	Fatal(c context.Context, msg ...interface{})

	Panic(c context.Context, msg ...interface{})
}
