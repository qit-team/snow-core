package utils

import (
	"testing"
)

func TestEncode62(t *testing.T) {
	num := int64(1122)
	s := Encode62(num)
	num1 := Decode62(s)
	if num != num1 {
		t.Error("it is not reversible")
		return
	}
}
