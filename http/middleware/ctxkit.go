package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/qit-team/snow-core/http/ctxkit"
	"github.com/qit-team/snow-core/log/logger"
)

func GenContextKit(c *gin.Context) {
	ctxkit.SetClientId(c, c.ClientIP())
	ctxkit.SetServerId(c, c.Request.RemoteAddr)
	ctxkit.SetHost(c, c.Request.Host)
	traceId := c.GetHeader("X-TRACE-ID")
	if traceId != "" {
		ctxkit.SetTraceId(c, traceId)
	} else {
		_, err := ctxkit.GenerateTraceId(c)
		if err != nil {
			logger.Error(c, "===GenContextKit.setTraceIdError===", err)
		}
	}
}
