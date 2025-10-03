package gmvc

import (
	"context"
	"io"
	"net/url"
)

type (

	// GmvcContext is the main context of gmvc. It will be used in gmvc everywhere.
	// It contains the HTTP request and response object.
	GmvcContext interface {
		context.Context

		// HttpRequest is the Raw HTTP request object.
		HttpRequest() HttpRequest

		// HttpResponse is the Raw HTTP response object.
		HttpResponse() HttpResponse

		// ActionMeta returns the current target Action's metadata.
		ActionMeta() *ActionMeta

		// Action returns the current target Action object.
		Action() any

		SetActionMeta(meta *ActionMeta)
		SetAction(action any)

		// GetCtx returns the value associated with this context for key, or nil if no
		// value is associated with key. Successive calls to GetCtx with the same key
		// returns the same result.
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

	// HttpRequest represents the general HTTP request object.
	HttpRequest interface {
		// Header info.

		// Method returns the HTTP request method.
		// Ref to: [net/http.MethodPost]
		Method() string

		// Host returns the host component of the request URL.
		Host() string

		// ContentLength returns the length of the request body.
		ContentLength() int

		// ContentType returns the Content-Type header of the request.
		ContentType() string

		// URL returns the request URL.
		URL() *url.URL

		// Header returns the request header.
		Header() Header

		// GetQuery returns the query parameter value for the named key.
		// If the key does not exist, it returns ("", false).
		GetQuery(key string) (string, bool)

		// GetPostForm returns the post form parameter value for the named key.
		// If the key does not exist, it returns ("", false).
		GetPostForm(key string) (string, bool)

		// GetForm returns the form parameter value for the named key.
		// If the key does not exist, it returns ("", false).
		GetForm(key string) (string, bool)

		// GetPathParam returns the path parameter value for the named key.
		// If the key does not exist, it returns ("", false).
		GetPathParam(key string) (string, bool)

		// VisitAllPostForm visits all post form parameters.
		VisitAllPostForm(func(key, value string))

		// VisitAllQuery visits all query parameters.
		VisitAllQuery(func(key, value string))

		// Body returns the request body.
		// FIXME: not sure if it is a good idea to return []byte. or maybe io.Reader is better.
		Body() []byte
	}

	// HttpResponse represents the general HTTP response object.
	HttpResponse interface {

		// FIXME: remove HTML method
		HTML(status int, body string, model any)

		// Status sets the HTTP response status code.
		Status(code int)

		// Header returns the response header.
		Header() Header

		// SetHeader sets the HTTP response header.
		SetHeader(key, value string)

		// Body sets the HTTP response body.
		Body(io.Reader)
	}

	// Header interface represents the HTTP header.
	Header interface {
		// Get returns the value associated with the given key.
		// If the key does not exist, it returns ("", false).
		Get(key string) (string, bool)

		// Gets returns the values associated with the given key.
		// If the key does not exist, it returns (nil, false).
		Gets(key string) ([]string, bool)

		// VisitAll visits all header key-value pairs.
		VisitAll(func(k, v []byte))
	}
)

// Initializer is invoked after the parameters are parsed.
type Initializer interface {

	// Init is invoked after the parameters are parsed.
	// If init returns an error, the request will be aborted and the error will be returned to the client.
	// Any error during gmvc runtime will be catched by [HandleError].
	Init() error
}

// Action is the main handler of gmvc. It will be invoked after the parameters are parsed and initialized.
type Action interface {
	Go() (interface{}, error)
}

// Validator is the validator function.
// It will be invoked after the parameters are resolved.
// If the validator returns an error, the request will be aborted and the error will be returned to the client.
// Any error during gmvc runtime will be catched by [HandleError].
type Validator func(ctx GmvcContext, fieldMeta *ParamMeta, value interface{}) error

// Resolver is the custom resolver function.
// It will be invoked after the parameters are parsed from the protocol.
// If the resolver returns an error, the request will be aborted and the error will be returned to the client.
// Any error during gmvc runtime will be catched by [HandleError].
type Resolver func(ctx GmvcContext, fieldMeta *ParamMeta, origin string) (interface{}, error)

// HandleError is the global error handler.
// If any error occurs during gmvc runtime, it will be catched by this handler.
// This handler convert the error to a response.
type HandleError func(ctx GmvcContext, err error) interface{}

// HandlerFunc defines the general handler function.
type HandlerFunc func(ctx GmvcContext)

// RecoverFunc is invoked when gmvc catches a panic.
// Any error during gmvc runtime will be catched by [HandleError].
type RecoverFunc func(ctx GmvcContext, info interface{}) (interface{}, error)
