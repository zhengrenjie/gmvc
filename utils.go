package gmvc

import (
	"reflect"
)

// 判断一个值是否是零值
func isZeroValue(value interface{}) bool {
	return value == nil || reflect.ValueOf(value).IsZero()
}

func hasSourceTag(target int, flag int) bool {
	return (target & flag) == flag
}
