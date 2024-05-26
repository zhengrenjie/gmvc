package gmvc

import (
	"errors"
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
	/* error definition */
	errParseParameter = errors.New("Parsing parameter error")

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

	convertStrings = []StringConvert[any]{
		convertAnyRet[string](convertString),
		convertAnyRet[*string](convertPointRet(convertString)),
		convertAnyRet[[]string](convertSliceRet(convertString)),
		convertAnyRet[[]*string](convertSliceRet(convertPointRet(convertString))),
	}

	convertBools = []StringConvert[any]{
		convertAnyRet[bool](convertBool),
		convertAnyRet[*bool](convertPointRet(convertBool)),
		convertAnyRet[[]bool](convertSliceRet(convertBool)),
		convertAnyRet[[]*bool](convertSliceRet(convertPointRet(convertBool))),
	}

	convertInts = []StringConvert[any]{
		convertAnyRet[int](convertInt),
		convertAnyRet[*int](convertPointRet(convertInt)),
		convertAnyRet[[]int](convertSliceRet(convertInt)),
		convertAnyRet[[]*int](convertSliceRet(convertPointRet(convertInt))),
	}

	convertInt8s = []StringConvert[any]{
		convertAnyRet[int8](convertInt8),
		convertAnyRet[*int8](convertPointRet(convertInt8)),
		convertAnyRet[[]int8](convertSliceRet(convertInt8)),
		convertAnyRet[[]*int8](convertSliceRet(convertPointRet(convertInt8))),
	}

	convertInt16s = []StringConvert[any]{
		convertAnyRet[int16](convertInt16),
		convertAnyRet[*int16](convertPointRet(convertInt16)),
		convertAnyRet[[]int16](convertSliceRet(convertInt16)),
		convertAnyRet[[]*int16](convertSliceRet(convertPointRet(convertInt16))),
	}

	convertInt32s = []StringConvert[any]{
		convertAnyRet[int32](convertInt32),
		convertAnyRet[*int32](convertPointRet(convertInt32)),
		convertAnyRet[[]int32](convertSliceRet(convertInt32)),
		convertAnyRet[[]*int32](convertSliceRet(convertPointRet(convertInt32))),
	}

	convertInt64s = []StringConvert[any]{
		convertAnyRet[int64](convertInt64),
		convertAnyRet[*int64](convertPointRet(convertInt64)),
		convertAnyRet[[]int64](convertSliceRet(convertInt64)),
		convertAnyRet[[]*int64](convertSliceRet(convertPointRet(convertInt64))),
	}

	convertUints = []StringConvert[any]{
		convertAnyRet[uint](convertUint),
		convertAnyRet[*uint](convertPointRet(convertUint)),
		convertAnyRet[[]uint](convertSliceRet(convertUint)),
		convertAnyRet[[]*uint](convertSliceRet(convertPointRet(convertUint))),
	}

	convertUint8s = []StringConvert[any]{
		convertAnyRet[uint8](convertUint8),
		convertAnyRet[*uint8](convertPointRet(convertUint8)),
		convertAnyRet[[]uint8](convertSliceRet(convertUint8)),
		convertAnyRet[[]*uint8](convertSliceRet(convertPointRet(convertUint8))),
	}

	convertUint16s = []StringConvert[any]{
		convertAnyRet[uint16](convertUint16),
		convertAnyRet[*uint16](convertPointRet(convertUint16)),
		convertAnyRet[[]uint16](convertSliceRet(convertUint16)),
		convertAnyRet[[]*uint16](convertSliceRet(convertPointRet(convertUint16))),
	}

	convertUint32s = []StringConvert[any]{
		convertAnyRet[uint32](convertUint32),
		convertAnyRet[*uint32](convertPointRet(convertUint32)),
		convertAnyRet[[]uint32](convertSliceRet(convertUint32)),
		convertAnyRet[[]*uint32](convertSliceRet(convertPointRet(convertUint32))),
	}

	convertUint64s = []StringConvert[any]{
		convertAnyRet[uint64](convertUint64),
		convertAnyRet[*uint64](convertPointRet(convertUint64)),
		convertAnyRet[[]uint64](convertSliceRet(convertUint64)),
		convertAnyRet[[]*uint64](convertSliceRet(convertPointRet(convertUint64))),
	}

	convertFloat32s = []StringConvert[any]{
		convertAnyRet[float32](convertFloat32),
		convertAnyRet[*float32](convertPointRet(convertFloat32)),
		convertAnyRet[[]float32](convertSliceRet(convertFloat32)),
		convertAnyRet[[]*float32](convertSliceRet(convertPointRet(convertFloat32))),
	}

	convertFloat64s = []StringConvert[any]{
		convertAnyRet[float64](convertFloat64),
		convertAnyRet[*float64](convertPointRet(convertFloat64)),
		convertAnyRet[[]float64](convertSliceRet(convertFloat64)),
		convertAnyRet[[]*float64](convertSliceRet(convertPointRet(convertFloat64))),
	}

	convertMap = map[reflect.Kind][]StringConvert[any]{
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
