package container

import "testing"

func TestContainer_SetSingleton(t *testing.T) {
	App.SetSingleton("di1", "1")
	App.SetSingleton("di2", 2)
	a1 := App.GetSingleton("di1")
	if a1 != "1" {
		t.Error("not same")
		return
	}

	a3 := App.GetSingleton("di3")
	if a3 != nil {
		t.Error("not same")
		return
	}
}
