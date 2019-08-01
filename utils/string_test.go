package utils

import (
	"testing"
)

func TestSubstr(t *testing.T) {
	str := "1234567890"
	s := Substr(str, 1, 2)
	if s != "23" {
		t.Error("substr failed")
		return
	}

	s = Substr(str, -2, 1)
	if s != "8" {
		t.Error("substr failed")
		return
	}

	s = Substr(str, -1, 0)
	if s != "" {
		t.Error("substr failed")
		return
	}
}

func TestJoin(t *testing.T) {
	s := Join("1", "2", "a")
	if s != "12a" {
		t.Error("join failed")
		return
	}
}
