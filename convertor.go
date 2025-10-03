package gmvc

import (
	"reflect"
	"strconv"
	"strings"
)

const (
	sliceSplit = ","
)

// StringConvert converts string to type T.
type StringConvert[T any] func(string) (T, error)

var (
	_ StringConvert[string] = convertString
	_ StringConvert[bool]   = convertBool

	_ StringConvert[int]   = convertInt
	_ StringConvert[int8]  = convertInt8
	_ StringConvert[int16] = convertInt16
	_ StringConvert[int32] = convertInt32
	_ StringConvert[int64] = convertInt64

	_ StringConvert[uint]   = convertUint
	_ StringConvert[uint8]  = convertUint8
	_ StringConvert[uint16] = convertUint16
	_ StringConvert[uint32] = convertUint32
	_ StringConvert[uint64] = convertUint64

	_ StringConvert[float32] = convertFloat32
	_ StringConvert[float64] = convertFloat64

	convertStrings = [4]StringConvert[any]{
		convertAnyRet(convertString),
		convertAnyRet(convertPointRet(convertString)),
		convertAnyRet(convertSliceRet(convertString)),
		convertAnyRet(convertSliceRet(convertPointRet(convertString))),
	}

	convertBools = [4]StringConvert[any]{
		convertAnyRet(convertBool),
		convertAnyRet(convertPointRet(convertBool)),
		convertAnyRet(convertSliceRet(convertBool)),
		convertAnyRet(convertSliceRet(convertPointRet(convertBool))),
	}

	convertInts = [4]StringConvert[any]{
		convertAnyRet(convertInt),
		convertAnyRet(convertPointRet(convertInt)),
		convertAnyRet(convertSliceRet(convertInt)),
		convertAnyRet(convertSliceRet(convertPointRet(convertInt))),
	}

	convertInt8s = [4]StringConvert[any]{
		convertAnyRet(convertInt8),
		convertAnyRet(convertPointRet(convertInt8)),
		convertAnyRet(convertSliceRet(convertInt8)),
		convertAnyRet(convertSliceRet(convertPointRet(convertInt8))),
	}

	convertInt16s = [4]StringConvert[any]{
		convertAnyRet(convertInt16),
		convertAnyRet(convertPointRet(convertInt16)),
		convertAnyRet(convertSliceRet(convertInt16)),
		convertAnyRet(convertSliceRet(convertPointRet(convertInt16))),
	}

	convertInt32s = [4]StringConvert[any]{
		convertAnyRet(convertInt32),
		convertAnyRet(convertPointRet(convertInt32)),
		convertAnyRet(convertSliceRet(convertInt32)),
		convertAnyRet(convertSliceRet(convertPointRet(convertInt32))),
	}

	convertInt64s = [4]StringConvert[any]{
		convertAnyRet(convertInt64),
		convertAnyRet(convertPointRet(convertInt64)),
		convertAnyRet(convertSliceRet(convertInt64)),
		convertAnyRet(convertSliceRet(convertPointRet(convertInt64))),
	}

	convertUints = [4]StringConvert[any]{
		convertAnyRet(convertUint),
		convertAnyRet(convertPointRet(convertUint)),
		convertAnyRet(convertSliceRet(convertUint)),
		convertAnyRet(convertSliceRet(convertPointRet(convertUint))),
	}

	convertUint8s = [4]StringConvert[any]{
		convertAnyRet(convertUint8),
		convertAnyRet(convertPointRet(convertUint8)),
		convertAnyRet(convertSliceRet(convertUint8)),
		convertAnyRet(convertSliceRet(convertPointRet(convertUint8))),
	}

	convertUint16s = [4]StringConvert[any]{
		convertAnyRet(convertUint16),
		convertAnyRet(convertPointRet(convertUint16)),
		convertAnyRet(convertSliceRet(convertUint16)),
		convertAnyRet(convertSliceRet(convertPointRet(convertUint16))),
	}

	convertUint32s = [4]StringConvert[any]{
		convertAnyRet(convertUint32),
		convertAnyRet(convertPointRet(convertUint32)),
		convertAnyRet(convertSliceRet(convertUint32)),
		convertAnyRet(convertSliceRet(convertPointRet(convertUint32))),
	}

	convertUint64s = [4]StringConvert[any]{
		convertAnyRet(convertUint64),
		convertAnyRet(convertPointRet(convertUint64)),
		convertAnyRet(convertSliceRet(convertUint64)),
		convertAnyRet(convertSliceRet(convertPointRet(convertUint64))),
	}

	convertFloat32s = [4]StringConvert[any]{
		convertAnyRet(convertFloat32),
		convertAnyRet(convertPointRet(convertFloat32)),
		convertAnyRet(convertSliceRet(convertFloat32)),
		convertAnyRet(convertSliceRet(convertPointRet(convertFloat32))),
	}

	convertFloat64s = [4]StringConvert[any]{
		convertAnyRet(convertFloat64),
		convertAnyRet(convertPointRet(convertFloat64)),
		convertAnyRet(convertSliceRet(convertFloat64)),
		convertAnyRet(convertSliceRet(convertPointRet(convertFloat64))),
	}

	convertMap = map[reflect.Kind][4]StringConvert[any]{
		reflect.Bool:    convertBools,
		reflect.String:  convertStrings,
		reflect.Int:     convertInts,
		reflect.Int8:    convertInt8s,
		reflect.Int16:   convertInt16s,
		reflect.Int32:   convertInt32s,
		reflect.Int64:   convertInt64s,
		reflect.Uint:    convertUints,
		reflect.Uint8:   convertUint8s,
		reflect.Uint16:  convertUint16s,
		reflect.Uint32:  convertUint32s,
		reflect.Uint64:  convertUint64s,
		reflect.Float32: convertFloat32s,
		reflect.Float64: convertFloat64s,
	}
)

func Convert(origin string, target reflect.Type) (interface{}, error) {
	isarray := target.Kind() == reflect.Slice
	if isarray {
		target = target.Elem()
	}

	ispoint := target.Kind() == reflect.Pointer
	if ispoint {
		target = target.Elem()
	}

	if !isarray && !ispoint {
		return convertMap[target.Kind()][0](origin)
	}

	if !isarray && ispoint {
		return convertMap[target.Kind()][1](origin)
	}

	if isarray && !ispoint {
		return convertMap[target.Kind()][2](origin)
	}

	return convertMap[target.Kind()][3](origin)
}

func convertSliceRet[T any](converter StringConvert[T]) StringConvert[[]T] {
	return func(s string) ([]T, error) {
		strs := strings.Split(s, sliceSplit)

		ret := make([]T, 0, len(strs))
		for _, str := range strs {
			ret0, err := converter(str)
			if err != nil {
				return nil, err
			}

			ret = append(ret, ret0)
		}

		return ret, nil
	}
}

func convertPointRet[T any](converter StringConvert[T]) StringConvert[*T] {
	return func(s string) (*T, error) {
		ret, err := converter(s)
		if err != nil {
			return nil, err
		}

		return &ret, nil
	}
}

func convertAnyRet[T any](converter StringConvert[T]) StringConvert[any] {
	return func(s string) (interface{}, error) {
		ret, err := converter(s)
		if err != nil {
			return nil, err
		}

		return ret, nil
	}
}

func convertString(s string) (string, error) {
	return s, nil
}

func convertBool(s string) (bool, error) {
	if s == "false" || s == "False" || s == "FALSE" {
		return false, nil
	}

	if s == "true" || s == "True" || s == "TRUE" {
		return true, nil
	}

	return false, nil // FIXME: return error
}

func convertInt(s string) (int, error) {
	value, err := strconv.ParseInt(s, 0, 0)
	if err != nil {
		return 0, err
	}

	return int(value), nil
}

func convertInt8(s string) (int8, error) {
	value, err := strconv.ParseInt(s, 0, 8)
	if err != nil {
		return 0, err
	}

	return int8(value), nil
}

func convertInt16(s string) (int16, error) {
	value, err := strconv.ParseInt(s, 0, 16)
	if err != nil {
		return 0, err
	}

	return int16(value), nil
}

func convertInt32(s string) (int32, error) {
	value, err := strconv.ParseInt(s, 0, 32)
	if err != nil {
		return 0, err
	}

	return int32(value), nil
}

func convertInt64(s string) (int64, error) {
	value, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func convertUint(s string) (uint, error) {
	value, err := strconv.ParseUint(s, 0, 0)
	if err != nil {
		return 0, err
	}

	return uint(value), nil
}

func convertUint8(s string) (uint8, error) {
	value, err := strconv.ParseUint(s, 0, 8)
	if err != nil {
		return 0, err
	}

	return uint8(value), nil
}

func convertUint16(s string) (uint16, error) {
	value, err := strconv.ParseUint(s, 0, 16)
	if err != nil {
		return 0, err
	}

	return uint16(value), nil
}

func convertUint32(s string) (uint32, error) {
	value, err := strconv.ParseUint(s, 0, 32)
	if err != nil {
		return 0, err
	}

	return uint32(value), nil
}

func convertUint64(s string) (uint64, error) {
	value, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func convertFloat32(s string) (float32, error) {
	value, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}

	return float32(value), nil
}

func convertFloat64(s string) (float64, error) {
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}
