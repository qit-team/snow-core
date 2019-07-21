package middleware

import (
	"github.com/qit-team/snow-core/log/accesslogger"
	"github.com/gin-gonic/gin"
	"github.com/hetiansu5/accesslog"
	"time"
)

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		//忽略HEAD探针的日志
		if c.Request.Method != "HEAD" {
			AccessLogFunc(accesslogger.GetAccessLogger())(c)
		}
	}
}

// AccessLogFunc 用于记录 http access log
func AccessLogFunc(accessLogger *accesslog.AccessLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		receivedAt := time.Now()
		originalWriter := c.Writer
		proxyWriter := newResponseWriter(c.Writer)
		c.Writer = proxyWriter.(gin.ResponseWriter)
		// Process request
		if c != nil {
			c.Next()
		}
		accessLogger.Log(proxyWriter, c.Request, receivedAt, time.Since(receivedAt))
		c.Writer = originalWriter
	}
}

type ResponseWriter struct {
	gin.ResponseWriter
	fbt time.Time
}

func (rw *ResponseWriter) FirstByteTime() time.Time {
	return rw.fbt
}

func (rw *ResponseWriter) WriteHeaderNow() {
	rw.ResponseWriter.WriteHeaderNow()
	if rw.fbt.IsZero() {
		rw.fbt = time.Now()
	}
}

func (rw *ResponseWriter) Write(data []byte) (n int, err error) {
	rw.WriteHeaderNow()
	return rw.ResponseWriter.Write(data)
}

func (rw *ResponseWriter) WriteString(s string) (n int, err error) {
	rw.WriteHeaderNow()
	return rw.ResponseWriter.WriteString(s)
}

func newResponseWriter(writer gin.ResponseWriter) accesslog.ResponseWriter {
	return &ResponseWriter{ResponseWriter: writer}
}
