package gmvc

import (
	"net/url"
)

type HttpRequest interface {
	Method() string
	Host() string
	GetQuery(key string) (string, bool)
	VisitAllQuery(func(key, value string))
	GetPostForm(key string) (string, bool)
	VisitAllPostForm(func(key, value string))
	FormValue(key string) (string, bool)
	URL() *url.URL
	Header() Header
	Body() []byte
	ContentLength() int
}

type Header interface {
	Get(key string) string
	VisitAll(func(k, v []byte))
	Set(key, v string)
	Header() []byte
}
