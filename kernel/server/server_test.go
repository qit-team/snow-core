package server

import "testing"

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
