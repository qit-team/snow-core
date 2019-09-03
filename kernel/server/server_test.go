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

func TestStop(t *testing.T) {
	go func() {
		WaitStop()
	}()

	//time.Sleep(1)
	go func() {
		Stop()
	}()
}

func TestCloseService(t *testing.T) {
	CloseService()
}

func TestHandleUserCmd(t *testing.T) {
	err := HandleUserCmd("cmd", "../../.env_pid")
	if err == nil {
		t.Error("unknown cmd error")
		return
	}

	err = HandleUserCmd("stop", "../../.env_pid")
	// process already finished
	if err == nil {
		t.Error("stop cmd error")
		return
	}

	err = HandleUserCmd("restart", "../../.env_pid")
	// process already finished
	if err == nil {
		t.Error("restart cmd error")
		return
	}

	// todo construct more cases, for exampla the process is running
}
