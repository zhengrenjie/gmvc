package gmvc_hertz

import (
	"net/url"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/zhengrenjie/gmvc"
)

var _ gmvc.HttpRequest = (*hertzReqAdapter)(nil)

func AcquireHertzReqAdapter(hertzCtx *app.RequestContext) gmvc.HttpRequest {
	req := &hertzReqAdapter{}
	req.hertzCtx = hertzCtx
	req.hertzReq = &hertzCtx.Request
	return req
}

type hertzReqAdapter struct {
	hertzCtx *app.RequestContext
	hertzReq *protocol.Request

	method string
	url    *url.URL
}

// Host implements gmvc.HttpRequest.
func (adapter *hertzReqAdapter) Host() string {
	return string(adapter.hertzReq.Host())
}

func (adapter *hertzReqAdapter) FormValue(key string) (string, bool) {
	if value, ok := adapter.hertzCtx.GetQuery(key); ok {
		return value, ok
	}

	return adapter.hertzCtx.GetPostForm(key)
}

// GetPostForm implements gmvc.HttpRequest.
func (adapter *hertzReqAdapter) GetPostForm(key string) (string, bool) {
	return adapter.hertzCtx.GetPostForm(key)
}

// GetQuery implements gmvc.HttpRequest.
func (adapter *hertzReqAdapter) GetQuery(key string) (string, bool) {
	return adapter.hertzCtx.GetQuery(key)
}

func (adapter *hertzReqAdapter) VisitAllQuery(f func(key, value string)) {
	adapter.hertzCtx.VisitAllQueryArgs(func(key, value []byte) {
		f(string(key), string(value))
	})
}

func (adapter *hertzReqAdapter) VisitAllPostForm(f func(key, value string)) {
	adapter.hertzCtx.VisitAllPostArgs(func(key, value []byte) {
		f(string(key), string(value))
	})
}

// Close implements gmvc.HttpRequest.
func (adapter *hertzReqAdapter) Close() error {
	return nil
}

/*- 实现HttpRequest -*/

func (adapter *hertzReqAdapter) Method() string {
	if adapter.method == "" {
		adapter.method = string(adapter.hertzReq.Method())
	}

	return adapter.method
}

func (adapter *hertzReqAdapter) URL() *url.URL {
	if adapter.url == nil {
		originUri := string(adapter.hertzReq.URI().FullURI())
		adapter.url, _ = url.ParseRequestURI(originUri)
	}

	return adapter.url
}

func (adapter *hertzReqAdapter) Header() gmvc.Header {
	return &adapter.hertzReq.Header
}

func (adapter *hertzReqAdapter) Body() []byte {
	return adapter.hertzReq.Body()
}

func (adapter *hertzReqAdapter) ContentLength() int {
	return adapter.hertzReq.Header.ContentLength()
}
