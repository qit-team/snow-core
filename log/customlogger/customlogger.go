package customlogger

import (
	"context"

	"github.com/qit-team/snow-core/log/logger"
	"github.com/qit-team/snow-core/log/utils"
)

// CustomLogger 自定义日志组件
type CustomLogger struct {
	FileName string // 可以扩展其它字段
}

func (cl *CustomLogger) Trace(c context.Context, msg ...interface{}) {
	withFields, newMsg := utils.SplitMsg(msg)
	data := utils.FormatLog(c, withFields...)
	logger.GetLoggerWithFileName(cl.FileName).WithFields(data).Trace(newMsg...)
}

func (cl *CustomLogger) Debug(c context.Context, msg ...interface{}) {
	withFields, newMsg := utils.SplitMsg(msg)
	data := utils.FormatLog(c, withFields...)
	logger.GetLoggerWithFileName(cl.FileName).WithFields(data).Debug(newMsg...)
}

func (cl *CustomLogger) Info(c context.Context, msg ...interface{}) {
	withFields, newMsg := utils.SplitMsg(msg)
	data := utils.FormatLog(c, withFields...)
	logger.GetLoggerWithFileName(cl.FileName).WithFields(data).Info(newMsg...)
}

func (cl *CustomLogger) Warn(c context.Context, msg ...interface{}) {
	withFields, newMsg := utils.SplitMsg(msg)
	data := utils.FormatLog(c, withFields...)
	logger.GetLoggerWithFileName(cl.FileName).WithFields(data).Warn(newMsg...)
}

func (cl *CustomLogger) Error(c context.Context, msg ...interface{}) {
	withFields, newMsg := utils.SplitMsg(msg)
	data := utils.FormatLog(c, withFields...)
	logger.GetLoggerWithFileName(cl.FileName).WithFields(data).Error(newMsg...)
}

func (cl *CustomLogger) Fatal(c context.Context, msg ...interface{}) {
	withFields, newMsg := utils.SplitMsg(msg)
	data := utils.FormatLog(c, withFields...)
	logger.GetLoggerWithFileName(cl.FileName).WithFields(data).Fatal(newMsg...)
}

func (cl *CustomLogger) Panic(c context.Context, msg ...interface{}) {
	withFields, newMsg := utils.SplitMsg(msg)
	data := utils.FormatLog(c, withFields...)
	logger.GetLoggerWithFileName(cl.FileName).WithFields(data).Fatal(newMsg...)
}
