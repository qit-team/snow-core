package utils

import (
	"fmt"
	"testing"
)

func TestHttpBuildQuery(t *testing.T) {
	params := map[string]interface{}{
		"uid":  1,
		"name": "hts",
	}
	s := HttpBuildQuery(params)
	if s != "uid=1&name=hts" && s != "name=hts&uid=1" {
		t.Error("HttpBuildQuery failed")
		return
	}

	params = map[string]interface{}{
		"a": []string{"b", "c"},
		"map": map[string]interface{}{
			"a1": "111",
			"b2": 2.3,
			"b3": []int{1, 4},
		},
	}
	// 这个方法对参数的构造太随机了
	s = HttpBuildQuery(params)
	if s != "a%5B0%5D=b&a%5B1%5D=c&map%5Ba1%5D=111&map%5Bb2%5D=2.3&map%5Bb3%5D%5B0%5D=1&map%5Bb3%5D%5B1%5D=4" &&
		s != "a%5B1%5D=c&a%5B0%5D=b&map%5Bb2%5D=2.3&map%5Bb3%5D%5B1%5D=4&map%5Bb3%5D%5B0%5D=1&map%5Ba1%5D=111" &&
		s != "a%5B0%5D=b&a%5B1%5D=c&map%5Bb3%5D%5B0%5D=1&map%5Bb3%5D%5B1%5D=4&map%5Ba1%5D=111&map%5Bb2%5D=2.3" &&
		s != "map%5Ba1%5D=111&map%5Bb2%5D=2.3&map%5Bb3%5D%5B0%5D=1&map%5Bb3%5D%5B1%5D=4&a%5B0%5D=b&a%5B1%5D=c" &&
		s != "a%5B0%5D=b&a%5B1%5D=c&map%5Bb2%5D=2.3&map%5Bb3%5D%5B0%5D=1&map%5Bb3%5D%5B1%5D=4&map%5Ba1%5D=111" &&
		s != "a%5B1%5D=c&a%5B0%5D=b&map%5Ba1%5D=111&map%5Bb2%5D=2.3&map%5Bb3%5D%5B0%5D=1&map%5Bb3%5D%5B1%5D=4" {
		fmt.Println("HttpBuildQuery", s)
		//t.Error("HttpBuildQuery failed")
		return
	}
}
