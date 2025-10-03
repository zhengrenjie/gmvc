package gmvc

import (
	"reflect"
)

// ActionMeta action's structured information.
type ActionMeta struct {

	// handler的名称
	handlerName string

	// 结构体类型
	handlerType reflect.Type

	// 结构体field数量
	fieldNum int

	// FieldMetaList
	fieldList []*ParamMeta
}

func (meta ActionMeta) GetName() string {
	return meta.handlerName
}

func (meta ActionMeta) GetFieldMeta() []*ParamMeta {
	return meta.fieldList
}

func (meta ActionMeta) GetAutowireInstances() map[string]any {
	autowire := make(map[string]any, 0)
	for _, fieldMeta := range meta.fieldList {
		if fieldMeta.autowire != "" && fieldMeta.value != nil {
			autowire[fieldMeta.autowire] = fieldMeta.value
		}
	}

	return autowire
}

// ParamMeta handler Field的结构化信息
type ParamMeta struct {
	// the origin value, used in autowire
	value interface{}

	actionMeta *ActionMeta

	// Field原始信息
	fieldType reflect.StructField

	// Tag原始信息
	tagInfo reflect.StructTag

	// Field的名字，用于从Query、Body、Header中找值
	fieldName string

	// 是否需要递归解析，如果有继承结构体的情况下，可能需要递归解析。
	isRecursive bool

	// 递归之后的结构体信息。
	handlerMeta *ActionMeta

	// Validators
	checkers []Validator

	// Resolver
	resolver Resolver

	// Field来源，Query,Body,Header
	source Src

	// 有没有设置default值
	hasDefault bool

	// default值
	def string

	autowire string
}

func (meta ParamMeta) GetActionMeta() *ActionMeta {
	return meta.actionMeta
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

func (meta ParamMeta) GetRecursive() *ActionMeta {
	return meta.handlerMeta
}

func (meta ParamMeta) GetAutowire() string {
	return meta.autowire
}
