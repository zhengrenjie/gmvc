package gmvc_hertz

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/zhengrenjie/gmvc"
)

var (
	silence                  = struct{}{}
	_       gmvc.GmvcContext = (*HertzContext)(nil)
)

func AcquireHertzContext(ctx context.Context, hertzCtx *app.RequestContext) *HertzContext {
	ret := &HertzContext{
		c:        ctx,
		ctx:      hertzCtx,
		paramSet: make(map[string]interface{}, 0),

		requset:  AcquireHertzReqAdapter(hertzCtx),
		response: AcquireHertzRespAdapter(hertzCtx),
	}

	return ret
}

type HertzContext struct {
	c        context.Context
	ctx      *app.RequestContext
	paramSet map[string]interface{}

	requset  *hertzReqAdapter
	response *hertzRespAdapter

	action     any
	actionMeta *gmvc.ActionMeta
}

// SetAction implements gmvc.GmvcContext.
func (h *HertzContext) SetAction(action any) {
	h.action = action
}

// SetActionMeta implements gmvc.GmvcContext.
func (h *HertzContext) SetActionMeta(meta *gmvc.ActionMeta) {
	h.actionMeta = meta
}

// Action implements gmvc.GmvcContext.
func (h *HertzContext) Action() any {
	return h.action
}

// ActionMeta implements gmvc.GmvcContext.
func (h *HertzContext) ActionMeta() *gmvc.ActionMeta {
	return h.actionMeta
}

// HttpRequest implements gmvc.GmvcContext.
func (h *HertzContext) HttpRequest() gmvc.HttpRequest {
	return h.requset
}

// HttpResponse implements gmvc.GmvcContext.
func (h *HertzContext) HttpResponse() gmvc.HttpResponse {
	return h.response
}

func (h *HertzContext) Deadline() (deadline time.Time, ok bool) {
	return h.c.Deadline()
}

func (h *HertzContext) Done() <-chan struct{} {
	return h.c.Done()
}

func (h *HertzContext) Err() error {
	return h.c.Err()
}

func (h *HertzContext) Value(key interface{}) interface{} {
	return h.c.Value(key)
}

func (h *HertzContext) GetCtx(key string) (interface{}, bool) {
	return h.ctx.Get(key)
}

func (h *HertzContext) GetString(key string) string {
	rawVal, ok := h.GetCtx(key)
	if !ok {
		return ""
	}

	value, ok := rawVal.(string)
	if !ok {
		return ""
	}

	return value
}

func (h *HertzContext) GetInt(key string) int {
	rawVal, ok := h.GetCtx(key)
	if !ok {
		return 0
	}

	value, ok := rawVal.(int)
	if !ok {
		return 0
	}

	return value
}

func (h *HertzContext) GetHeader(key string) (string, bool) {
	rawHeader := string(h.ctx.GetHeader(key))
	if rawHeader == "" {
		return "", false
	}

	return rawHeader, true
}

func (h *HertzContext) GetHeaders(key string) ([]string, bool) {
	rawHeader := string(h.ctx.GetHeader(key))
	if rawHeader == "" {
		return nil, false
	}

	return []string{rawHeader}, true
}

func (h *HertzContext) GetQuery(key string) (string, bool) {
	return h.ctx.GetQuery(key)
}

func (h *HertzContext) GetForm(key string) (string, bool) {
	return h.ctx.GetPostForm(key)
}

func (h *HertzContext) GetPathParam(key string) (string, bool) {
	value := h.ctx.Param(key)
	if value == "" {
		return "", false
	}

	return value, true
}

func (h *HertzContext) GetRawData() ([]byte, error) {
	return h.ctx.GetRawData(), nil
}

func (h *HertzContext) GetContentType() string {
	return string(h.ctx.ContentType())
}

func (h *HertzContext) HasParam(name string) bool {
	_, ok := h.paramSet[name]
	return ok
}

func (h *HertzContext) Report(name string) {
	h.paramSet[name] = silence
}

func (h *HertzContext) GetStatus() int {
	return h.ctx.Response.StatusCode()
}

func (h *HertzContext) Abort() {
	h.ctx.Abort()
}

func (h *HertzContext) Set(key string, value interface{}) {
	h.ctx.Set(key, value)
}

func (h *HertzContext) GetEntity() interface{} {
	return h.ctx
}

func (h *HertzContext) Status(code int) {
	h.ctx.Status(code)
}

func (h *HertzContext) Header(key, value string) {
	h.ctx.Header(key, value)
}

func (h *HertzContext) HTML(name string, code int, obj interface{}) {
	panic("hertz not support templete render yet")
}

/* implements gmvc.HttpRequest, gmvc.Header, gmvc.HttpResponse */

var _ gmvc.HttpRequest = (*hertzReqAdapter)(nil)
var _ gmvc.Header = (*hertzReqHeaderAdapter)(nil)
var _ gmvc.HttpResponse = (*hertzRespAdapter)(nil)
var _ gmvc.Header = (*hertzRespHeaderAdapter)(nil)

func AcquireHertzReqAdapter(hertzCtx *app.RequestContext) *hertzReqAdapter {
	req := &hertzReqAdapter{}
	req.hertzCtx = hertzCtx
	req.hertzReq = &hertzCtx.Request
	req.header = &hertzReqHeaderAdapter{
		header: &hertzCtx.Request.Header,
	}
	return req
}

func AcquireHertzRespAdapter(hertzCtx *app.RequestContext) *hertzRespAdapter {
	req := &hertzRespAdapter{}
	req.hertzCtx = hertzCtx
	req.header = &hertzRespHeaderAdapter{
		header: &hertzCtx.Response.Header,
	}
	return req
}

type (
	hertzReqAdapter struct {
		hertzCtx *app.RequestContext
		hertzReq *protocol.Request

		method string
		url    *url.URL
		header *hertzReqHeaderAdapter
	}

	hertzRespAdapter struct {
		hertzCtx *app.RequestContext
		header   *hertzRespHeaderAdapter
	}

	hertzReqHeaderAdapter struct {
		header *protocol.RequestHeader
	}

	hertzRespHeaderAdapter struct {
		header *protocol.ResponseHeader
	}
)

// Get implements gmvc.Header.
func (h *hertzRespHeaderAdapter) Get(key string) (string, bool) {
	return h.header.Get(key), true
}

// Gets implements gmvc.Header.
func (h *hertzRespHeaderAdapter) Gets(key string) ([]string, bool) {
	return h.header.GetAll(key), true
}

// VisitAll implements gmvc.Header.
func (h *hertzRespHeaderAdapter) VisitAll(f func(k []byte, v []byte)) {
	h.header.VisitAll(f)
}

// Body implements gmvc.HttpResponse.
func (h *hertzRespAdapter) Body(out io.Reader) {
	h.hertzCtx.SetBodyStream(out, -1)
}

// HTML implements gmvc.HttpResponse.
func (h *hertzRespAdapter) HTML(status int, body string, model any) {
	h.hertzCtx.HTML(status, body, model)
}

// Header implements gmvc.HttpResponse.
func (h *hertzRespAdapter) Header() gmvc.Header {
	return h.header
}

// SetHeader implements gmvc.HttpResponse.
func (h *hertzRespAdapter) SetHeader(key string, value string) {
	h.hertzCtx.Header(key, value)
}

// Status implements gmvc.HttpResponse.
func (h *hertzRespAdapter) Status(code int) {
	h.hertzCtx.Status(code)
}

// Get implements gmvc.Header.
func (h *hertzReqHeaderAdapter) Get(key string) (string, bool) {
	return h.header.Get(key), true
}

// Gets implements gmvc.Header.
func (h *hertzReqHeaderAdapter) Gets(key string) ([]string, bool) {
	return h.header.GetAll(key), true
}

// VisitAll implements gmvc.Header.
func (h *hertzReqHeaderAdapter) VisitAll(f func(k []byte, v []byte)) {
	h.header.VisitAll(f)
}

// ContentType implements gmvc.HttpRequest.
func (adapter *hertzReqAdapter) ContentType() string {
	return string(adapter.hertzReq.Header.ContentType())
}

// GetForm implements gmvc.HttpRequest.
func (adapter *hertzReqAdapter) GetForm(key string) (string, bool) {
	v := adapter.hertzCtx.FormValue(key)
	if v == nil {
		return "", false
	}

	return string(v), true
}

// GetPathParam implements gmvc.HttpRequest.
func (adapter *hertzReqAdapter) GetPathParam(key string) (string, bool) {
	return adapter.hertzCtx.Params.Get(key)
}

// Host implements gmvc.HttpRequest.
func (adapter *hertzReqAdapter) Host() string {
	return string(adapter.hertzReq.Host())
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
	return adapter.header
}

func (adapter *hertzReqAdapter) Body() []byte {
	return adapter.hertzReq.Body()
}

func (adapter *hertzReqAdapter) ContentLength() int {
	return adapter.hertzReq.Header.ContentLength()
}
