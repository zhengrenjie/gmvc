package gmvc

import (
	"reflect"
	"strconv"
	"strings"
)

// Convert 将string转换成需要的type
type Convert func(string, reflect.Type) (interface{}, error)

var (
	convertMap map[reflect.Kind]Convert
	sliceSplit = ","
)

func init() {
	convertMap = make(map[reflect.Kind]Convert)

	convertMap[reflect.Int] = func(originValue string, targetType reflect.Type) (interface{}, error) {
		value, err := strconv.ParseInt(originValue, 10, 0)
		if err != nil {
			return nil, err
		}

		return int(value), nil
	}

	convertMap[reflect.Int8] = func(originValue string, targetType reflect.Type) (interface{}, error) {
		value, err := strconv.ParseInt(originValue, 10, 8)
		if err != nil {
			return nil, err
		}

		return int8(value), nil
	}

	convertMap[reflect.Int16] = func(originValue string, targetType reflect.Type) (interface{}, error) {
		value, err := strconv.ParseInt(originValue, 10, 16)
		if err != nil {
			return nil, err
		}

		return int16(value), nil
	}

	convertMap[reflect.Int32] = func(originValue string, targetType reflect.Type) (interface{}, error) {
		value, err := strconv.ParseInt(originValue, 10, 32)
		if err != nil {
			return nil, err
		}

		return int32(value), nil
	}

	convertMap[reflect.Int64] = func(originValue string, targetType reflect.Type) (interface{}, error) {
		value, err := strconv.ParseInt(originValue, 10, 64)
		if err != nil {
			return nil, err
		}

		return int64(value), nil
	}

	convertMap[reflect.Uint] = func(originValue string, targetType reflect.Type) (interface{}, error) {
		value, err := strconv.ParseUint(originValue, 10, 0)
		if err != nil {
			return nil, err
		}

		return uint(value), nil
	}

	convertMap[reflect.Uint8] = func(originValue string, targetType reflect.Type) (interface{}, error) {
		value, err := strconv.ParseUint(originValue, 10, 8)
		if err != nil {
			return nil, err
		}

		return uint8(value), nil
	}

	convertMap[reflect.Uint16] = func(originValue string, targetType reflect.Type) (interface{}, error) {
		value, err := strconv.ParseUint(originValue, 10, 16)
		if err != nil {
			return nil, err
		}

		return uint16(value), nil
	}

	convertMap[reflect.Uint32] = func(originValue string, targetType reflect.Type) (interface{}, error) {
		value, err := strconv.ParseUint(originValue, 10, 32)
		if err != nil {
			return nil, err
		}

		return uint32(value), nil
	}

	convertMap[reflect.Uint64] = func(originValue string, targetType reflect.Type) (interface{}, error) {
		value, err := strconv.ParseUint(originValue, 10, 64)
		if err != nil {
			return nil, err
		}

		return uint64(value), nil
	}

	convertMap[reflect.Float32] = func(originValue string, targetType reflect.Type) (interface{}, error) {
		value, err := strconv.ParseFloat(originValue, 32)
		if err != nil {
			return nil, err
		}

		return float32(value), nil
	}

	convertMap[reflect.Float64] = func(originValue string, targetType reflect.Type) (interface{}, error) {
		value, err := strconv.ParseFloat(originValue, 64)
		if err != nil {
			return nil, err
		}

		return float64(value), nil
	}

	convertMap[reflect.String] = func(originValue string, targetType reflect.Type) (interface{}, error) {
		return originValue, nil
	}

	convertMap[reflect.Ptr] = func(originValue string, targetType reflect.Type) (interface{}, error) {
		targetType = targetType.Elem()
		convert, ok := convertMap[targetType.Kind()]
		if ok {
			value, err := convert(originValue, targetType)
			if err != nil {
				return nil, err
			}

			switch targetType.Kind() {
			case reflect.Int8:
				ret, _ := (value.(int8))
				return &ret, nil
			case reflect.Int16:
				ret, _ := (value.(int16))
				return &ret, nil
			case reflect.Int32:
				ret, _ := (value.(int32))
				return &ret, nil
			case reflect.Int64:
				ret, _ := (value.(int64))
				return &ret, nil
			case reflect.Int:
				ret, _ := (value.(int))
				return &ret, nil
			case reflect.String:
				ret, _ := (value.(string))
				return &ret, nil
			case reflect.Uint8:
				ret, _ := (value.(uint8))
				return &ret, nil
			case reflect.Uint16:
				ret, _ := (value.(uint16))
				return &ret, nil
			case reflect.Uint32:
				ret, _ := (value.(uint32))
				return &ret, nil
			case reflect.Uint64:
				ret, _ := (value.(uint64))
				return &ret, nil
			case reflect.Uint:
				ret, _ := (value.(uint))
				return &ret, nil
			case reflect.Float32:
				ret, _ := (value.(float32))
				return &ret, nil
			case reflect.Float64:
				ret, _ := (value.(float64))
				return &ret, nil
			default:
				return nil, nil
			}
		}

		return nil, nil
	}

	convertMap[reflect.Slice] = func(originValue string, targetType reflect.Type) (interface{}, error) {
		componentType := targetType.Elem()
		strSplit := strings.Split(originValue, sliceSplit)

		var sliceValue interface{}

		switch componentType.Kind() {
		case reflect.Int8:
			sliceValue = make([]int8, 0, len(strSplit))
		case reflect.Int16:
			sliceValue = make([]int16, 0, len(strSplit))
		case reflect.Int32:
			sliceValue = make([]int32, 0, len(strSplit))
		case reflect.Int64:
			sliceValue = make([]int64, 0, len(strSplit))
		case reflect.Int:
			sliceValue = make([]int, 0, len(strSplit))
		case reflect.String:
			sliceValue = make([]string, 0, len(strSplit))
		case reflect.Uint8:
			sliceValue = make([]uint8, 0, len(strSplit))
		case reflect.Uint16:
			sliceValue = make([]uint16, 0, len(strSplit))
		case reflect.Uint32:
			sliceValue = make([]uint32, 0, len(strSplit))
		case reflect.Uint64:
			sliceValue = make([]uint64, 0, len(strSplit))
		case reflect.Uint:
			sliceValue = make([]uint, 0, len(strSplit))
		case reflect.Float32:
			sliceValue = make([]float32, 0, len(strSplit))
		case reflect.Float64:
			sliceValue = make([]float64, 0, len(strSplit))
		default:
			sliceValue = make([]interface{}, 0, len(strSplit))
		}

		for _, item := range strSplit {
			converted, err := convertMap[componentType.Kind()](item, nil)
			if err != nil {
				return nil, err
			}

			switch componentType.Kind() {
			case reflect.Int8:
				sliceValue = append(sliceValue.([]int8), converted.(int8))
			case reflect.Int16:
				sliceValue = append(sliceValue.([]int16), converted.(int16))
			case reflect.Int32:
				sliceValue = append(sliceValue.([]int32), converted.(int32))
			case reflect.Int64:
				sliceValue = append(sliceValue.([]int64), converted.(int64))
			case reflect.Int:
				sliceValue = append(sliceValue.([]int), converted.(int))
			case reflect.String:
				sliceValue = append(sliceValue.([]string), converted.(string))
			case reflect.Uint8:
				sliceValue = append(sliceValue.([]uint8), converted.(uint8))
			case reflect.Uint16:
				sliceValue = append(sliceValue.([]uint16), converted.(uint16))
			case reflect.Uint32:
				sliceValue = append(sliceValue.([]uint32), converted.(uint32))
			case reflect.Uint64:
				sliceValue = append(sliceValue.([]uint64), converted.(uint64))
			case reflect.Uint:
				sliceValue = append(sliceValue.([]uint), converted.(uint))
			case reflect.Float32:
				sliceValue = append(sliceValue.([]float32), converted.(float32))
			case reflect.Float64:
				sliceValue = append(sliceValue.([]float64), converted.(float64))
			default:
				sliceValue = append(sliceValue.([]interface{}), converted)
			}
		}

		return sliceValue, nil
	}
}
