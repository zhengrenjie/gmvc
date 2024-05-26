package gmvc

import (
	"time"

	"github.com/gin-gonic/gin"
)

var (
	silence = struct{}{}
)

type GinContext struct {
	ctx      *gin.Context
	paramSet map[string]interface{}
}

func (g *GinContext) Deadline() (deadline time.Time, ok bool) {
	return g.ctx.Deadline()
}

func (g *GinContext) Done() <-chan struct{} {
	return g.ctx.Done()
}

func (g *GinContext) Err() error {
	return g.ctx.Err()
}

func (g *GinContext) Value(key interface{}) interface{} {
	return g.ctx.Value(key)
}

func (g *GinContext) GetCtx(key string) (interface{}, bool) {
	return g.ctx.Get(key)
}

func (g *GinContext) GetString(key string) string {
	rawVal, ok := g.GetCtx(key)
	if !ok {
		return ""
	}

	value, ok := rawVal.(string)
	if !ok {
		return ""
	}

	return value
}

func (g *GinContext) GetInt(key string) int {
	rawVal, ok := g.GetCtx(key)
	if !ok {
		return 0
	}

	value, ok := rawVal.(int)
	if !ok {
		return 0
	}

	return value
}

func (g *GinContext) GetHeader(key string) (string, bool) {
	header, ok := g.ctx.Request.Header[key]
	if ok {
		return header[0], true
	}

	return "", false
}

func (g *GinContext) GetHeaders(key string) ([]string, bool) {
	header, ok := g.ctx.Request.Header[key]
	if ok {
		return header, true
	}

	return nil, false
}

func (g *GinContext) GetQuery(key string) (string, bool) {
	return g.ctx.GetQuery(key)
}

func (g *GinContext) GetForm(key string) (string, bool) {
	return g.ctx.GetPostForm(key)
}

func (g *GinContext) GetPathParam(key string) (string, bool) {
	value := g.ctx.Param(key)
	if value == "" {
		return "", false
	}

	return value, true
}

func (g *GinContext) GetRawData() ([]byte, error) {
	return g.ctx.GetRawData()
}

func (g *GinContext) GetContentType() string {
	return g.ctx.ContentType()
}

func (g *GinContext) HasParam(name string) bool {
	_, ok := g.paramSet[name]
	return ok
}

func (g *GinContext) Report(name string) {
	g.paramSet[name] = silence
}

func (g *GinContext) GetStatus() int {
	return g.ctx.Writer.Status()
}

func (g *GinContext) Abort() {
	g.ctx.Abort()
}

func (g *GinContext) Set(key string, value interface{}) {
	g.ctx.Set(key, value)
}

func (g *GinContext) GetEntity() interface{} {
	return g.ctx
}

func (g *GinContext) Status(code int) {
	g.ctx.Status(code)
}

func (g *GinContext) Header(key, value string) {
	g.ctx.Header(key, value)
}

func (g *GinContext) JSON(code int, obj interface{}) {
	g.ctx.JSON(code, obj)
}

func (g *GinContext) HTML(name string, code int, obj interface{}) {
	g.HTML(name, code, obj)
}

func (g *GinContext) String(code int, value string) {
	g.ctx.String(code, value)
}

var _ HttpContext = &GinContext{}

type GinAdapter struct {
	*DefaultDispatcher
}

var GinInstance = &GinAdapter{
	DefaultDispatcher: DefaultInstance,
}

func (Gin *GinAdapter) Wrap(h Action) gin.HandlerFunc {
	handler := Gin.InnerWrap(h)
	return func(c *gin.Context) {
		adapter := &GinContext{
			ctx:      c,
			paramSet: make(map[string]interface{}, 0),
		}
		handler(adapter)
	}
}
