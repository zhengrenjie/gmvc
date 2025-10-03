package gmvc_hertz

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/zhengrenjie/gmvc"
)

var (
	silence                  = struct{}{}
	_       gmvc.GmvcContext = (*HertzContext)(nil)
)

type HertzContext struct {
	c        context.Context
	ctx      *app.RequestContext
	paramSet map[string]interface{}

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
	return AcquireHertzReqAdapter(h.ctx)
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

func (h *HertzContext) JSON(code int, obj interface{}) {
	h.ctx.JSON(code, obj)
}

func (h *HertzContext) HTML(name string, code int, obj interface{}) {
	panic("hertz not support templete render yet")
}

func (h *HertzContext) String(code int, value string) {
	h.ctx.String(code, value)
}
