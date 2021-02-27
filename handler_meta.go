package gmvc

import (
	"reflect"
	"strings"
	"unicode"
)

// HandlerMeta handler的结构化信息
type HandlerMeta struct {

	// handler的名称
	handlerName string

	// 结构体类型
	handlerType reflect.Type

	// 结构体field数量
	fieldNum int

	// FieldMetaList
	fieldList []*ParamMeta
}

func (meta HandlerMeta) GetName() string {
	return meta.handlerName
}

func (meta HandlerMeta) GetFieldMeta() []*ParamMeta {
	return meta.fieldList
}

// ParamMeta handler Field的结构化信息
type ParamMeta struct {
	// Field原始信息
	fieldType reflect.StructField

	// Tag原始信息
	tagInfo reflect.StructTag

	// Field的名字，用于从Query、Body、Header中找值
	fieldName string

	// 是否需要递归解析，如果有继承结构体的情况下，可能需要递归解析。
	isRecursive bool

	// 递归之后的结构体信息。
	handlerMeta *HandlerMeta

	// Validators
	checkers []Validator

	// Resolver
	resolver Resolver

	// Field来源，Query,Body,Header
	source int

	// 有没有设置default值
	hasDefault bool

	// default值
	def string
}

func (meta ParamMeta) GetName() string {
	return meta.fieldName
}

func (meta ParamMeta) GetType() reflect.Type {
	return meta.fieldType.Type
}

func (meta ParamMeta) GetTag() reflect.StructTag {
	return meta.tagInfo
}

func (meta ParamMeta) GetResolver() Resolver {
	return meta.resolver
}

func (meta ParamMeta) GetValidators() []Validator {
	return meta.checkers
}

func (meta ParamMeta) GetDefault() string {
	return meta.def
}

func (meta ParamMeta) GetRecursive() *HandlerMeta {
	return meta.handlerMeta
}

// HandlerParser 把handler parse成内部的结构体
type HandlerParser struct {
	dispatcher *DefaultDispatcher
}

func (parser *HandlerParser) introspect(struct0 reflect.Type) *HandlerMeta {
	numField := struct0.NumField()
	structMeta := &HandlerMeta{
		handlerName: struct0.Name(),
		handlerType: struct0,
		fieldNum:    numField,
		fieldList:   make([]*ParamMeta, 0),
	}

	for i := 0; i < numField; i++ {
		field := struct0.Field(i)
		tagInfo := field.Tag
		fieldMeta := &ParamMeta{
			fieldType: field,
			tagInfo:   tagInfo,
			fieldName: field.Name,
		}

		// xParams解析
		xParams := strings.Split(tagInfo.Get(XParam), XSplit)
		for _, value := range xParams {
			if value == "" {
				continue
			}
			switch value {
			case XRecursive:
				// 如果需要递归解析，则递归下去。
				fieldMeta.isRecursive = true
				fieldMeta.handlerMeta = parser.introspect(field.Type)
			case XQuery:
				fieldMeta.source |= QuerySrc
			case XBody:
				fieldMeta.source |= BodySrc
			case XHeader:
				fieldMeta.source |= HeaderSrc
				// header 首字母自动大写
				name := fieldMeta.fieldName
				for i, v := range name {
					fieldMeta.fieldName = string(unicode.ToUpper(v)) + name[i+1:]
					break
				}
			case XPath:
				fieldMeta.source |= PathSrc
			case XCtx:
				fieldMeta.source |= CtxSrc
			case XAuto:
				fieldMeta.source |= AutoSrc
			default:
				fieldMeta.fieldName = value
			}
		}

		// xValidator解析
		validatorStr, ok := tagInfo.Lookup(XValidator)
		if ok {
			xValidator := strings.Split(validatorStr, XSplit)
			checkers := make([]Validator, len(xValidator))
			fieldMeta.checkers = checkers
			for index, value := range xValidator {
				if value == "" {
					continue
				}
				checker := parser.dispatcher.checkerMap[value]
				if checker == nil {
					continue
				}
				checkers[index] = checker
			}
		}

		// Resolver解析
		resolverStr, ok := tagInfo.Lookup(XResolver)
		if ok {
			resolver, ok := parser.dispatcher.resolverMap[resolverStr]
			if ok {
				fieldMeta.resolver = resolver
			}
		}

		// xDefault解析
		defaultStr, ok := tagInfo.Lookup(XDefault)
		if ok {
			fieldMeta.hasDefault = true
			fieldMeta.def = defaultStr
		}

		structMeta.fieldList = append(structMeta.fieldList, fieldMeta)
	}

	return structMeta
}
