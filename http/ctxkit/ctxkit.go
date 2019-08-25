package ctxkit

import (
	"github.com/gin-gonic/gin"
	"context"
	"github.com/qit-team/snow-core/utils"
	"crypto/md5"
	"fmt"
	"strings"
	"errors"
)

const (
	TraceId  = "x-trace-id"
	ClientIp = "x-cip"
	ServerIp = "x-sip"
	HOST     = "x-host"
)

func SetTraceId(ctx context.Context, value string) {
	if ctxGin, ok := ctx.(*gin.Context); ok {
		SetGinTraceId(ctxGin, value)
	}
}
func SetGinTraceId(ctx *gin.Context, value string) {
	ctx.Set(TraceId, value)
}

func GetTraceId(ctx context.Context) string {
	s, _ := ctx.Value(TraceId).(string)
	return s
}

func SetClientId(ctx context.Context, value string) {
	if ctxGin, ok := ctx.(*gin.Context); ok {
		SetGinClientId(ctxGin, value)
	}
}

func SetGinClientId(ctx *gin.Context, value string) {
	ctx.Set(ClientIp, value)
}

func GetClientId(ctx context.Context) string {
	s, _ := ctx.Value(ClientIp).(string)
	return s
}

// param to change
func SetServerId(ctx *gin.Context, value string) {
	ctx.Set(ServerIp, value)
}

func GetServerId(ctx context.Context) string {
	s, _ := ctx.Value(ServerIp).(string)
	return s
}

// param to change
func SetHost(ctx *gin.Context, value string) {
	ctx.Set(HOST, value)
}

func GetHost(ctx context.Context) string {
	s, _ := ctx.Value(HOST).(string)
	return s
}

//var once sync.Once
func GenerateTraceId(ctx context.Context) (string, error) {
	randomId := utils.GenUUID()
	mdTemp := md5.Sum([]byte(randomId))
	mdCode := fmt.Sprintf("%x", mdTemp)
	mdStr := strings.ToUpper(mdCode)
    if len(mdStr) < 32 {
    	return "", errors.New("grenerateTraceIdError")
	}
	traceId := mdStr[0 : 8] + "-" + mdStr[8 : 12] + "-" + mdStr[12 : 16] + "-" + mdStr[16 : 20] + "-" + mdStr[20 : 32]

	SetTraceId(ctx, traceId)
	return traceId, nil
}