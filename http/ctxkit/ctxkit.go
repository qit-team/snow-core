package ctxkit

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qit-team/snow-core/utils"
	"strings"
)

const (
	TraceId  = "x-trace-id"
	ClientIp = "x-cip"
	ServerIp = "x-sip"
	HOST     = "x-host"
)

func SetTraceId(ctx context.Context, value string) context.Context {
	var newCtx context.Context
	if ctxGin, ok := ctx.(*gin.Context); ok {
		newCtx = SetGinTraceId(ctxGin.Request.Context(), value)
		ctxGin.Request = ctxGin.Request.WithContext(newCtx)
	} else {
		newCtx = SetGinTraceId(ctx, value)
	}
	return newCtx
}

func SetGinTraceId(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, TraceId, value)
}

func GetTraceId(ctx context.Context) string {
	if ctxGin, ok := ctx.(*gin.Context); ok {
		ctx = ctxGin.Request.Context()
	}
	s, _ := ctx.Value(TraceId).(string)
	return s
}

func SetClientId(ctx context.Context, value string) context.Context {
	var newCtx context.Context
	if ctxGin, ok := ctx.(*gin.Context); ok {
		newCtx = SetGinClientId(ctxGin.Request.Context(), value)
		ctxGin.Request = ctxGin.Request.WithContext(SetGinClientId(ctxGin.Request.Context(), value))
	} else {
		newCtx = SetGinClientId(ctx, value)
	}
	return newCtx
}

func SetGinClientId(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, ClientIp, value)
}

func GetClientId(ctx context.Context) string {
	if ctxGin, ok := ctx.(*gin.Context); ok {
		ctx = ctxGin.Request.Context()
	}
	s, _ := ctx.Value(ClientIp).(string)
	return s
}

func SetServerId(ctx context.Context, value string) context.Context {
	var newCtx context.Context
	if ctxGin, ok := ctx.(*gin.Context); ok {
		newCtx = SetGinServerId(ctxGin.Request.Context(), value)
		ctxGin.Request = ctxGin.Request.WithContext(newCtx)
	} else {
		newCtx = SetGinServerId(ctx, value)
	}
	return newCtx
}

func SetGinServerId(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, ServerIp, value)
}

func GetServerId(ctx context.Context) string {
	if ctxGin, ok := ctx.(*gin.Context); ok {
		ctx = ctxGin.Request.Context()
	}
	s, _ := ctx.Value(ServerIp).(string)
	return s
}

// param to change
func SetHost(ctx context.Context, value string) context.Context {
	var newCtx context.Context
	if ctxGin, ok := ctx.(*gin.Context); ok {
		newCtx = SetGinHost(ctxGin.Request.Context(), value)
		ctxGin.Request = ctxGin.Request.WithContext(newCtx)
	} else {
		newCtx = SetGinHost(ctx, value)
	}
	return newCtx
}

func SetGinHost(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, ServerIp, value)
}

func GetHost(ctx context.Context) string {
	if ctxGin, ok := ctx.(*gin.Context); ok {
		ctx = ctxGin.Request.Context()
	}
	s, _ := ctx.Value(HOST).(string)
	return s
}

//var once sync.Once
func GenerateTraceId(ctx context.Context) (string, context.Context) {
	randomId := utils.GenUUID()
	mdTemp := md5.Sum([]byte(randomId))
	mdCode := fmt.Sprintf("%x", mdTemp)
	mdStr := strings.ToUpper(mdCode)

	var traceId = mdStr
	if len(mdStr) >= 32 {
		traceId = mdStr[0:8] + "-" + mdStr[8:12] + "-" + mdStr[12:16] + "-" + mdStr[16:20] + "-" + mdStr[20:32]
	}
	return traceId, SetTraceId(ctx, traceId)
}
