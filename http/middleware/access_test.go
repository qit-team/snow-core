package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hetiansu5/accesslog"
	"testing"
)

var handle, handleFunc gin.HandlerFunc
var accessLogger *accesslog.AccessLogger

var reponseWriter accesslog.ResponseWriter

func init() {
	handle = AccessLog()
	handleFunc = AccessLogFunc(accessLogger)
	var w gin.ResponseWriter
	reponseWriter = newResponseWriter(w)
}

func TestResponseWriter(t *testing.T) {
	ret := reponseWriter.FirstByteTime()
	fmt.Println("FirstByteTime", ret)
}
