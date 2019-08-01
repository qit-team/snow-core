package ctxkit

import (
	"testing"
	"github.com/gin-gonic/gin"
)

var c *gin.Context

func init() {
	c = &gin.Context{}
}

func TestGetClientId(t *testing.T) {
	v := "1"
	SetClientId(c, v)
	v1 := GetClientId(c)
	if v1 != v {
		t.Error("ClientId miss match")
		return
	}
}

func TestGetTraceId(t *testing.T) {
	v := "2"
	SetTraceId(c, v)
	v1 := GetTraceId(c)
	if v1 != v {
		t.Error("TraceId miss match")
		return
	}
}

func TestGetHost(t *testing.T) {
	v := "3"
	SetHost(c, v)
	v1 := GetHost(c)
	if v1 != v {
		t.Error("Host miss match")
		return
	}
}

func TestGetServerId(t *testing.T) {
	v := "4"
	SetServerId(c, v)
	v1 := GetServerId(c)
	if v1 != v {
		t.Error("ServerId miss match")
		return
	}
}
