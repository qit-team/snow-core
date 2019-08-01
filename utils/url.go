package utils

import (
	"fmt"
	"reflect"
	"strings"
	"net/url"
)

// MapToStringList 多层Map转字符串数组
func mapToStringList(params map[string]interface{}, parentNode string) (result []string) {
	for key := range params {
		nextParentNode := ""
		if len(parentNode) > 0 {
			nextParentNode = parentNode + "[" + key + "]"
		} else {
			nextParentNode = key
		}
		value := params[key]
		t := reflect.TypeOf(value)
		switch t.Kind() {
		case reflect.Map:
			tempResult := mapToStringList(value.(map[string]interface{}), nextParentNode)
			result = append(result, tempResult...)
			break
		case reflect.Slice:
			typeString := t.Elem().String()
			tmpVal := map[string]interface{}{}

			if typeString == "int" {
				for idx, subVal := range value.([]int) {
					tmpVal[fmt.Sprint(idx)] = subVal
				}
			} else if typeString == "string" {
				for idx, subVal := range value.([]string) {
					tmpVal[fmt.Sprint(idx)] = subVal
				}
			} else {
				for idx, subVal := range value.([]interface{}) {
					tmpVal[fmt.Sprint(idx)] = subVal
				}
			}
			tempResult := mapToStringList(tmpVal, nextParentNode)
			result = append(result, tempResult...)
			break
		default:
			result = append(result, url.QueryEscape(nextParentNode)+"="+url.QueryEscape(fmt.Sprint(value)))
		}
	}
	return
}

//HttpBuildQuery 生成Query参数
func HttpBuildQuery(params map[string]interface{}) (query string) {
	queryList := mapToStringList(params, "")
	query = strings.Join(queryList, "&")
	return
}
