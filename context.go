package gmvc

import (
	"context"
	"net/url"
)

// GmvcContext http的通用context
type (
	GmvcContext interface {
		context.Context

		// 获取原始HttpRequest
		HttpRequest() HttpRequest

		HttpResponse() HttpResponse

		// return action meta
		ActionMeta() *ActionMeta

		// return action object
		Action() any

		SetActionMeta(meta *ActionMeta)

		SetAction(action any)

		// 获取上下文参数
		GetCtx(key string) (interface{}, bool)

		// 是否含有某个参数
		HasParam(name string) bool

		// 参数报道
		Report(name string)

		// set context
		Set(key string, value interface{})

		// 获取真实的Context对象，例如 gin.Context, hertz.RequestContext
		GetEntity() interface{}
	}

	HttpRequest interface {
		Method() string
		Host() string
		ContentLength() int
		ContentType() string
		GetQuery(key string) (string, bool)
		GetPostForm(key string) (string, bool)
		GetForm(key string) (string, bool)
		GetPathParam(key string) (string, bool)
		VisitAllPostForm(func(key, value string))

		VisitAllQuery(func(key, value string))
		URL() *url.URL
		Header() Header
		Body() []byte
	}

	HttpResponse interface {
		Status(code int)
		Header(key, value string)
		Body() []byte
	}

	Header interface {
		Get(key string) (string, bool)
		Gets(key string) ([]string, bool)
		VisitAll(func(k, v []byte))
		Set(key, v string)
	}
)
