package utils

import "testing"

func TestGetCurrentTime(t *testing.T) {
	t1 := GetCurrentTime()
	t2 := GetCurrentMilliTime()
	if t1 != t2/1000 && t1+1 != t2/1000 {
		t.Error("time error")
		return
	}
}
