package gmvc

import (
	"reflect"
)

func (instance *DefaultDispatcher) resolve(c HTTPContext, meta *HandlerMeta) (interface{}, error) {
	// 实例化handler
	handlerValuePtr := reflect.New(meta.handlerType)

	// 解析每个结构体的值
	if err := instance.resolveAndCheckFieldValue(c, handlerValuePtr, meta); err != nil {
		return nil, err
	}

	return handlerValuePtr.Interface(), nil
}

func (instance *DefaultDispatcher) resolveAndCheckFieldValue(ctx HTTPContext, pvalue reflect.Value, meta *HandlerMeta) error {
	for i := 0; i < meta.fieldNum; i++ {
		fieldMeta := meta.fieldList[i]

		if fieldMeta.isRecursive {
			recursiveValueStr := reflect.New(fieldMeta.handlerMeta.handlerType)
			if resp := instance.resolveAndCheckFieldValue(ctx, recursiveValueStr, fieldMeta.handlerMeta); resp != nil {
				return resp
			}
			pvalue.Elem().Field(i).Set(recursiveValueStr.Elem())
		}

		var value interface{}

		if reflect.TypeOf(ctx).AssignableTo(fieldMeta.fieldType.Type) {
			value = ctx
		} else if reflect.TypeOf(ctx.GetEntity()).AssignableTo(fieldMeta.fieldType.Type) {
			value = ctx.GetEntity()
		} else {
			originValue, src, ok := instance.drawOutOriginValue(ctx, fieldMeta)

			/*
			 found parameter
			 1. report param found.
			 2. if param found in context, just assign to value.
			 3. if param has resolver, then use the return of resolver as the value.
			 4. auto type-convert.
			*/
			if ok {
				ctx.Report(fieldMeta.fieldName)

				if src == CtxSrc {
					value = originValue
				} else if fieldMeta.resolver != nil {
					var err error
					if value, err = fieldMeta.resolver(ctx, fieldMeta, originValue.(string)); err != nil {
						return err
					}
				} else {
					value = instance.resolveFieldValue(ctx, fieldMeta, originValue.(string), meta.handlerName)
				}
			}

			if err := instance.validateValue(ctx, value, fieldMeta, meta); err != nil {
				return err
			}
		}

		if value == nil {
			continue
		}

		valueType := reflect.TypeOf(value)
		if !valueType.AssignableTo(fieldMeta.fieldType.Type) {
			continue
		}

		switch valueType.Kind() {
		case reflect.Int:
			pvalue.Elem().Field(i).SetInt(int64(value.(int)))
		case reflect.Int8:
			pvalue.Elem().Field(i).SetInt(int64(value.(int8)))
		case reflect.Int16:
			pvalue.Elem().Field(i).SetInt(int64(value.(int16)))
		case reflect.Int32:
			pvalue.Elem().Field(i).SetInt(int64(value.(int32)))
		case reflect.Int64:
			pvalue.Elem().Field(i).SetInt(value.(int64))
		case reflect.String:
			pvalue.Elem().Field(i).SetString(value.(string))
		case reflect.Uint:
			pvalue.Elem().Field(i).SetUint(uint64(value.(uint)))
		case reflect.Uint8:
			pvalue.Elem().Field(i).SetUint(uint64(value.(uint8)))
		case reflect.Uint16:
			pvalue.Elem().Field(i).SetUint(uint64(value.(uint16)))
		case reflect.Uint32:
			pvalue.Elem().Field(i).SetUint(uint64(value.(uint32)))
		case reflect.Uint64:
			pvalue.Elem().Field(i).SetUint(value.(uint64))
		case reflect.Float32:
			pvalue.Elem().Field(i).SetFloat(float64(value.(float32)))
		case reflect.Float64:
			pvalue.Elem().Field(i).SetFloat(value.(float64))
		default:
			pvalue.Elem().Field(i).Set(reflect.ValueOf(value))
		}
	}

	return nil
}

func (instance *DefaultDispatcher) drawOutOriginValue(ctx HTTPContext, fieldMeta *ParamMeta) (originValue interface{}, src int, present bool) {
	if hasSourceTag(fieldMeta.source, HeaderSrc) {
		src = HeaderSrc
		originValue, present = ctx.GetHeader(fieldMeta.fieldName)
		if present {
			return
		}
	}

	if hasSourceTag(fieldMeta.source, QuerySrc) {
		src = QuerySrc
		originValue, present = ctx.GetQuery(fieldMeta.fieldName)
		if present {
			return
		}
	}

	if hasSourceTag(fieldMeta.source, PathSrc) {
		src = PathSrc
		originValue, present = ctx.GetPathParam(fieldMeta.fieldName)
		if present {
			return
		}
	}

	if hasSourceTag(fieldMeta.source, BodySrc) {
		src = BodySrc
		originValue, present = ctx.GetForm(fieldMeta.fieldName)
		if present {
			return
		}
	}

	if hasSourceTag(fieldMeta.source, CtxSrc) {
		src = CtxSrc
		originValue, present = ctx.GetCtx(fieldMeta.fieldName)
		if present {
			return
		}
	}

	if fieldMeta.hasDefault {
		src = DefaultSrc
		present = true
		originValue = fieldMeta.def
	}

	return
}

func (instance *DefaultDispatcher) validateValue(ctx HTTPContext, value interface{}, fieldMeta *ParamMeta, meta *HandlerMeta) error {
	if len(fieldMeta.checkers) > 0 {
		for _, checker := range fieldMeta.checkers {
			if checker == nil {
				continue
			}

			err := checker(ctx, fieldMeta, value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 参数类型转换
func (instance *DefaultDispatcher) resolveFieldValue(ctx HTTPContext, fieldMeta *ParamMeta, originValue string, handlerName string) interface{} {
	convert, ok := convertMap[fieldMeta.fieldType.Type.Kind()]
	if ok {
		value, err := convert(originValue, fieldMeta.fieldType.Type)
		if err != nil {
			return nil
		}

		return value
	}

	return nil
}
