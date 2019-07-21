package ctxkit

import (
	"github.com/gin-gonic/gin"
	"context"
)

const (
	TraceId  = "x-trace-id"
	ClientIp = "x-cip"
	ServerIp = "x-sip"
	HOST     = "x-host"
)

func SetTraceId(ctx *gin.Context, value string) {
	ctx.Set(TraceId, value)
}

func GetTraceId(ctx context.Context) string {
	s, _ := ctx.Value(TraceId).(string)
	return s
}

func SetClientId(ctx *gin.Context, value string) {
	ctx.Set(ClientIp, value)
}

func GetClientId(ctx context.Context) string {
	s, _ := ctx.Value(ClientIp).(string)
	return s
}

func SetServerId(ctx *gin.Context, value string) {
	ctx.Set(ServerIp, value)
}

func GetServerId(ctx context.Context) string {
	s, _ := ctx.Value(ServerIp).(string)
	return s
}

func SetHost(ctx *gin.Context, value string) {
	ctx.Set(HOST, value)
}

func GetHost(ctx context.Context) string {
	s, _ := ctx.Value(HOST).(string)
	return s
}
