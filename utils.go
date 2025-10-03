package gmvc

import (
	"reflect"
)

// 判断一个值是否是零值
func isZeroValue(value interface{}) bool {
	return value == nil || reflect.ValueOf(value).IsZero()
}

func hasSourceTag(target Src, flag Src) bool {
	return (target & flag) == flag
}
