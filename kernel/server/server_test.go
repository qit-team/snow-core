package server

import (
	"os"
	"syscall"
	"testing"
)

func TestGetDebug(t *testing.T) {
	debug := GetDebug()
	if debug != false {
		t.Error("debug status is error")
		return
	}
	SetDebug(true)
	debug = GetDebug()
	if debug != true {
		t.Error("debug status is error")
		return
	}
}

func TestSignel(t *testing.T) {
	// 这种函数只能通过单测跑一遍看是否有报错 没有返回数据orError类型可以判断
	RegisterSignal()

	go func() {
		var sigs = []os.Signal{
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
		}

		HandleSignal(sigs[0])
		HandleSignal(sigs[1])
		HandleSignal(sigs[2])
	}()
}

func TestPidFile(t *testing.T) {
	err := WritePidFile("../../.env_pid", 10001)
	if err != nil {
		t.Error("WritePidFile error")
		return
	}

	pid, err := ReadPidFile("../../.env_pid")
	if err != nil {
		t.Error("ReadPidFile error")
		return
	} else if pid != 10001 {
		t.Error("ReadPidFile error result not right")
		return
	}
}

// todo
// github.com/qit-team/snow-core/kernel/server/server.go:68:		WaitStop			0.0%
// github.com/qit-team/snow-core/kernel/server/server.go:73:		CloseService			0.0%
// github.com/qit-team/snow-core/kernel/server/server.go:111:		HandleUserCmd			0.0%
// github.com/qit-team/snow-core/kernel/server/server.go:139:		Stop				0.0%
