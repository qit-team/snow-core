package utils

import (
	"context"
	"os"

	"github.com/qit-team/snow-core/http/ctxkit"
	"github.com/sirupsen/logrus"
)

var hostname string

type withField struct {
	Key   string
	Value interface{}
}

//SplitMsg 将日志消息分裂
func SplitMsg(msg []interface{}) (withFields []*withField, newMsg []interface{}) {
	for _, v := range msg {
		switch v.(type) {
		case *withField:
			withFields = append(withFields, v.(*withField))
		case []*withField:
			// 如果是通过batchNewWithFields，需要做如下处理
			tempWithFieldsList := v.([]*withField)
			if len(tempWithFieldsList) != 0 {
				for _, tempWithField := range tempWithFieldsList {
					withFields = append(withFields, tempWithField)
				}
			}
		default:
			newMsg = append(newMsg, v)
		}
	}
	return
}

func FormatLog(c context.Context, args ...*withField) logrus.Fields {
	data := logrus.Fields{
		"host": GetHostName(),
	}

	if c != nil {
		traceId := ctxkit.GetTraceId(c)
		if traceId != "" {
			data["trace_id"] = traceId
		} else {
			traceId, _ := ctxkit.GenerateTraceId(c)
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

func GetHostName() string {
	if hostname == "" {
		hostname, _ = os.Hostname()
		if hostname == "" {
			hostname = "unknown"
		}
	}
	return hostname
}
