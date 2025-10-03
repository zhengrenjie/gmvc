package gmvc

import (
	"context"
)

// GmvcContext http的通用context
type GmvcContext interface {
	context.Context

	// return action meta
	ActionMeta() *ActionMeta

	// return action object
	Action() any

	SetActionMeta(meta *ActionMeta)

	SetAction(action any)

	// 获取原始HttpRequest
	HttpRequest() HttpRequest

	// 获取上下文参数
	GetCtx(key string) (interface{}, bool)

	// 上下文中取string
	GetString(key string) string

	// 上下文中取int
	GetInt(key string) int

	// 获取header
	GetHeader(key string) (string, bool)

	// 获取headers
	GetHeaders(key string) ([]string, bool)

	// 获取参数，Query
	GetQuery(key string) (string, bool)

	// 获取body中的form数据
	GetForm(key string) (string, bool)

	// 获取path的数据
	GetPathParam(key string) (string, bool)

	// 获取body的raw data
	GetRawData() ([]byte, error)

	// 获取ContentType
	GetContentType() string

	// 是否含有某个参数
	HasParam(name string) bool

	// 参数报道
	Report(name string)

	// 返回status code
	GetStatus() int

	// 中断
	Abort()

	// set context
	Set(key string, value interface{})

	// 获取真实的Context对象，例如 gin.Context, hertz.RequestContext
	GetEntity() interface{}

	// Response 状态码
	Status(code int)

	// response set header
	Header(key, value string)

	// json 返回
	JSON(code int, obj interface{})

	// html 返回
	HTML(name string, code int, obj interface{})

	// string 返回
	String(code int, value string)
}
