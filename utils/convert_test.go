package utils

import (
	"testing"
)

func TestMapStrInterface2MapStrStr(t *testing.T) {
	mp := map[string]interface{}{
		"a": 1,
		"b": "bb",
		"c": 3.2,
		"d": false,
	}
	mp1 := MapStrInterface2MapStrStr(mp)
	if mp1["a"] != "1" {
		t.Error("map a failed")
		return
	} else if mp1["b"] != "bb" {
		t.Error("map b failed")
		return
	} else if mp1["c"] != "3.2" {
		t.Error("map c failed")
		return
	} else if mp1["d"] != "false" {
		t.Error("map d failed")
		return
	}
}

func TestNum2Str(t *testing.T) {
	num := 2211
	s := Num2Str(num)
	if s != "2211" {
		t.Error("Num2Str failed")
		return
	}
}

func TestSliceStr2Interface(t *testing.T) {
	arr := []string{
		"a", "b",
	}
	arr2 := SliceStr2Interface(arr)
	if !(arr2[0] == "a" && arr2[1] == "b") {
		t.Error("SliceStr2Interface failed")
		return
	}
}
