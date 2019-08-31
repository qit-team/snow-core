package container

import (
	"fmt"
	"strings"
	"testing"
)

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

func TestContainer_demo(t *testing.T) {
	nameStr := App.String()

	if len(nameStr) == 0 {
		t.Error("String() empty")
		return
	}

	if strings.Index(nameStr, "di1") == -1 || strings.Index(nameStr, "di2") == -1 {
		t.Error("String ret error")
		return
	}
	fmt.Println("=======string ret:", nameStr)

	bool1 := App.isSingleton("di1")
	if !bool1 {
		t.Error("isSingleton Error")
		return
	}

	bool2 := App.isPrototype("snow-test")
	if bool2 {
		t.Error("isPrototype Error")
		return
	}

	strTest := App.injectName("snow-test,snow-test111")
	if strTest != "snow-test" || len(strTest) == 0 {
		t.Error("injectName Error")
		return
	}

	//App.Ensure("di1")

	App.SetPrototype("snow", factoryDemo)

	ret, err := App.GetPrototype("snow")
	if err != nil {
		fmt.Println("=======", err)
		t.Error("GetPrototype error")
		return
	}
	fmt.Println("GetPrototype ret:", ret)

	// after set prototype string again ,for cover more branch of if&else
	nameStr = App.String()

	if len(nameStr) == 0 {
		t.Error("String() empty")
		return
	}

	// for cover branch of exception return
	strTest = App.injectName("")
	if len(strTest) != 0 {
		t.Error("injectName Exception branch Error")
		return
	}

	bool1 = App.isSingleton("prototype")
	if bool1 {
		t.Error("isSingleton Exception branch Error")
		return
	}

	bool2 = App.isPrototype("prototype")
	if !bool2 {
		t.Error("isPrototype Exception branch Error")
		return
	}

}

func factoryDemo() (i interface{}, err error) {
	return
}
