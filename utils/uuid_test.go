package utils

import (
	"testing"
)

func TestGenUUID(t *testing.T) {
	s := GenUUID()
	if len(s) == 0 {
		t.Error("length of uuid is equal 0")
	}
}
