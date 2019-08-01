package alimns

import (
	"testing"
	"github.com/qit-team/snow-core/config"
)

func Test_getSingleton(t *testing.T) {
	c := getSingleton("", false)
	if c != nil {
		t.Error("client is not equal nil")
		return
	}
}

func TestProvider(t *testing.T) {
	err := Pr.Register("mns", config.MnsConfig{}, true)
	if err != nil {
		t.Error(err)
		return
	}

	arr := Pr.Provides()
	if !(len(arr) == 1 && arr[0] == "mns") {
		t.Errorf("Provides is not match. %v", arr)
		return
	}

	err = Pr.Register("mns1", config.MnsConfig{})
	if err != nil {
		t.Error(err)
		return
	}

	arr = Pr.Provides()
	if !(len(arr) == 2 && arr[1] == "mns1"|| arr[1] == "mns") {
		t.Errorf("Provides is not match. %v", arr)
		return
	}

	err = Pr.Close()
	if err != nil {
		t.Error(err)
		return
	}

	c := GetMns()
	if c != nil {
		t.Error("client is not equal nil")
		return
	}
}
