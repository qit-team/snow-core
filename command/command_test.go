package command

import (
	"testing"
	"fmt"
)

func TestNew(t *testing.T) {
	cmd := New()
	cmd.AddFunc("test", test)
	cmd.Execute("test")

	defer func() {
		if e := recover(); e == nil {
			t.Error("unknown name do not panic")
		}
	}()
	cmd.Execute("test1")
}

func test() {
	fmt.Println("run test")
}
