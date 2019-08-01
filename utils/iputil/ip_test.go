package iputil

import (
	"testing"
)

//内网IP
func TestGetInternalIp(t *testing.T) {
	str, err := GetInternalIp()
	if err != nil {
		t.Error(err)
		return
	} else if len(str) == 0 {
		t.Error("get internal ip failed")
		return
	}
}
