package logger

import (
	"github.com/sirupsen/logrus"
	"context"
	"os"
	"github.com/qit-team/snow-core/http/ctxkit"
)

var (
	hostname string
)

type withField struct {
	Key   string
	Value interface{}
}

//此结构的数据将会在挂靠到日志的一级键中体现
//demo: logger.Info(ctx, "curl", NewWithFiled("key1", "value1"), NewWithFiled("key2", "value2"), "msg1", "msg2")
func NewWithField(key string, value interface{}) *withField {
	return &withField{Key: key, Value: value}
}

//批量
func BatchNewWithField(data map[string]interface{}) (arr []*withField) {
	for k, v := range data {
		arr = append(arr, NewWithField(k, v))
	}
	return arr
}

func GetHostName() string {
	if hostname == "" {
		hostname, _ = os.Hostname()
		if hostname == "" {
			hostname = "unknown"
		}
	}
	return hostname
}

func formatLog(c context.Context, t string, args ...*withField) logrus.Fields {
	data := logrus.Fields{
		"type": t,
		"host": GetHostName(),
	}

	if c != nil {
		traceId := ctxkit.GetTraceId(c)
		if traceId != "" {
			data["trace_id"] = traceId
		}

		domain := ctxkit.GetHost(c)
		if domain != "" {
			data["domain"] = domain
		}

		sip := ctxkit.GetServerId(c)
		if sip != "" {
			data["sip"] = sip
		}

		cip := ctxkit.GetClientId(c)
		if cip != "" {
			data["cip"] = cip
		}
	}

	for _, field := range args {
		if _, ok := data[field.Key]; !ok {
			data[field.Key] = field.Value
		}
	}

	return data
}

func Trace(c context.Context, logType string, msg ...interface{}) {
	withFields, newMsg := splitMsg(msg)
	data := formatLog(c, logType, withFields...)
	GetLogger().WithFields(data).Trace(newMsg...)
}

func Debug(c context.Context, logType string, msg ...interface{}) {
	withFields, newMsg := splitMsg(msg)
	data := formatLog(c, logType, withFields...)
	GetLogger().WithFields(data).Debug(newMsg...)
}

func Info(c context.Context, logType string, msg ...interface{}) {
	withFields, newMsg := splitMsg(msg)
	data := formatLog(c, logType, withFields...)
	GetLogger().WithFields(data).Info(newMsg...)
}

func Warn(c context.Context, logType string, msg ...interface{}) {
	withFields, newMsg := splitMsg(msg)
	data := formatLog(c, logType, withFields...)
	GetLogger().WithFields(data).Warn(newMsg...)
}

func Error(c context.Context, logType string, msg ...interface{}) {
	withFields, newMsg := splitMsg(msg)
	data := formatLog(c, logType, withFields...)
	GetLogger().WithFields(data).Error(newMsg...)
}

func Fatal(c context.Context, logType string, msg ...interface{}) {
	withFields, newMsg := splitMsg(msg)
	data := formatLog(c, logType, withFields...)
	GetLogger().WithFields(data).Fatal(newMsg...)
}

func Panic(c context.Context, logType string, msg ...interface{}) {
	withFields, newMsg := splitMsg(msg)
	data := formatLog(c, logType, withFields...)
	GetLogger().WithFields(data).Panic(newMsg...)
}

//将日志消息分裂
func splitMsg(msg []interface{}) (withFields []*withField, newMsg []interface{}) {
	for _, v := range msg {
		switch v.(type) {
		case *withField:
			withFields = append(withFields, v.(*withField))
		default:
			newMsg = append(newMsg, v)
		}
	}
	return
}
