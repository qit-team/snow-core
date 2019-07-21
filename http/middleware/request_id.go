package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/qit-team/snow-core/utils"
)

func GenRequestId(c *gin.Context) {
	reqId := utils.GenUUID()
	c.Request.Header.Add("X-Request-Id", reqId)
	c.Header("X-Request-Id", reqId)
	c.Next()
}
