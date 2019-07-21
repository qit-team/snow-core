package helper

import (
	"errors"
)

func GetDiName(defaultName string, args ...string) string {
	var name string
	if len(args) > 0 {
		name = args[0]
	}
	if name == "" {
		return defaultName
	}
	return name
}

func TransformArgs(args ...interface{}) (diName string, lazy bool, err error) {
	if len(args) < 2 {
		err = errors.New("args is not enough")
		return
	}

	var ok bool
	diName, ok = args[0].(string)
	if !ok {
		err = errors.New("args[0] is not string")
		return
	}

	if len(args) > 2 {
		lazy, _ = args[2].(bool)
	}
	return
}

func MapToArray(mp map[string]interface{}) []string {
	arr := make([]string, len(mp))
	i := 0
	for k := range mp {
		arr[i] = k
		i++
	}
	return arr
}
