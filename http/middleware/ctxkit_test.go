package middleware

import (
	"github.com/gin-gonic/gin"
	"testing"
)

var c *gin.Context

func init() {
	c = &gin.Context{}
}

func Test_GenContextKit(t *testing.T) {
	//c.Header("X-Forwarded-For", "127.0.0.111")
	//c1 := gin.Context{}
	////fmt.Println("c.engine.ForwardedByClientIP", c.requestHeader("X-Forwarded-For"))
	//fmt.Println("=======111111")
	//GenContextKit(c1)
	//fmt.Println("========2222") // 校验traceId是否设置成功
	//traceId := ctxkit.GetTraceId(c)
	//if len(traceId) == 0 {
	//	t.Error("GenContextKit error")
	//	return
	//}
	//fmt.Println("GenContextKit traceId:", traceId)
}
