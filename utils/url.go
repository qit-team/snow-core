package utils

import (
	"net/url"
	"fmt"
)

//对map结构进行http_query_build
func HttpBuildQuery(params map[string]interface{}) string {
	v := url.Values{}
	for key, value := range params {
		v.Add(key, fmt.Sprint(value))
	}
	return v.Encode()
}
