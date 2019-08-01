package helper

import "testing"

func TestTransformArgs(t *testing.T) {
	_, _, err := TransformArgs("1")
	if err == nil {
		t.Error("length of args should be checked")
		return
	}

	_, _, err = TransformArgs(1, "", true)
	if err == nil {
		t.Error("args[0] should be string")
		return
	}

	diName, lazy, err := TransformArgs("1", "", true)
	if err != nil {
		t.Error(err)
		return
	} else if diName != "1" {
		t.Error("diName is not match")
		return
	} else if lazy != true {
		t.Error("lazy is not match")
		return
	}

}

func TestGetDiName(t *testing.T) {
	dn := "dn"
	a1 := GetDiName(dn)
	if a1 != dn {
		t.Error("must be default")
		return
	}

	a2 := GetDiName(dn, "22")
	if a2 != "22" {
		t.Error("must be args[0]")
		return
	}
}

func TestMapToArray(t *testing.T) {
	mp := map[string]interface{}{
		"a1": 1,
		"b2": "bbd",
	}
	arr := MapToArray(mp)
	if len(arr) != 2 {
		t.Error("length of array is not equal 2")
		return
	}

	if arr[0] == "a1" {
		if arr[1] != "b2" {
			t.Error("part result of array is error")
			return
		}
	} else if arr[0] == "b2" {
		if arr[1] != "a1" {
			t.Error("part result of array is error")
			return
		}
	} else {
		t.Error("result of array is error")
		return
	}
}
