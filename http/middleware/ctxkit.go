package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/qit-team/snow-core/http/ctxkit"
)

func GenContextKit(c *gin.Context) {
	ctxkit.SetClientId(c, c.ClientIP())
	ctxkit.SetServerId(c, c.Request.RemoteAddr)
	ctxkit.SetHost(c, c.Request.Host)
	traceId := c.GetHeader("X-Trace-Id")
	if traceId != "" {
		ctxkit.SetTraceId(c, traceId)
	}
}
