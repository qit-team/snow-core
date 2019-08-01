package utils

import (
	"testing"
)

func TestGetMd5Hash(t *testing.T) {
	s := GetMd5Hash("ss")
	if len(s) != 32 {
		t.Errorf("length of md5 string is not equal 16. %s", s)
		return
	}
}
