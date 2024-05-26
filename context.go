package gmvc

import (
	"context"
	"io"
	"net/url"
)

// HttpContext
type HttpContext interface {
	context.Context

	Request() IHttpRequest
	Response() IHttpResponse

	// GetCtx gets value from context.
	GetCtx(key string) (interface{}, bool)

	// SetCtx sets context value.
	SetCtx(key string, value interface{})

	JSON(code int, obj interface{})
}

type IHttpRequest interface {
	Method() string
	Host() string
	URL() *url.URL
	Header() IRequestHeader
	GetQuery(key string) (string, bool)
	GetForm(key string) (string, bool)
	GetPostForm(key string) (string, bool)
	GetPathParam(key string) (string, bool)
	Querys() *url.Values
	Body() io.ReadCloser
}

type IHttpResponse interface {
	Header() IResponseHeader
	SetResponse(ctx context.Context, resp *Response)
	SetStatusCode(code int)
}

type IHeader interface {
	VisitAll(func(k, v string))
	GetHeader(key string) string
	GetHeaders(key string) []string
}

type IRequestHeader interface {
	IHeader
}

type IResponseHeader interface {
	IHeader
	Add(k, v string)
}
