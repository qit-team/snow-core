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

//var once sync.Once
func GenerateTraceId(ctx *gin.Context) (string, error) {
	randomId := utils.GenUUID()
	mdTemp := md5.Sum([]byte(randomId))
	mdCode := fmt.Sprintf("%x", mdTemp)
	mdStr := strings.ToUpper(mdCode)
    if len(mdStr) < 32 {
    	return "", errors.New("grenerateTraceIdError")
	}
	traceId := mdStr[0 : 8] + "-" + mdStr[8 : 12] + "-" + mdStr[12 : 16] + "-" + mdStr[16 : 20] + "-" + mdStr[20 : 32]
	SetTraceId(ctx, traceId)

	////加并发控制
	//once.Do(func() {
	//	SetTraceId(ctx, traceId)
	//})
	return traceId, nil
}