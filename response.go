package gmvc

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

// RenderType 返回渲染类型
type RenderType int32

const (
	// JSON: content-type="application/json"
	JSON RenderType = 0

	// HTML: content-type="text/html"
	HTML RenderType = 1

	// String: content-type="text/plain"
	String RenderType = 2
)

// Response is the convenient struct to return HTTP response.
// By seting [RenderType], gmvc will use different responsor to return response.
type Response struct {
	StatusCode int
	Body       interface{}
	Render     RenderType
	Header     map[string]string
	Model      map[string]interface{}
}

// Responsor 负责返回，可能会有多种Render
type Responsor interface {
	Response(ctx GmvcContext, resp *Response)
}

// JSONResponsor Json实现
type JSONResponsor struct{}

// Response 返回的方法
func (r *JSONResponsor) Response(ctx GmvcContext, resp *Response) {
	setDefault(resp)
	setHeader(ctx, resp)

	b, err := json.Marshal(resp.Body)
	if err != nil {
		ctx.HttpResponse().Status(http.StatusInternalServerError)
		return
	}

	ctx.HttpResponse().SetHeader("Content-Type", "application/json")
	ctx.HttpResponse().Status(resp.StatusCode)
	ctx.HttpResponse().Body(bytes.NewReader(b))
}

// HTMLResponsor html实现
type HTMLResponsor struct{}

// Response 返回的方法
func (r *HTMLResponsor) Response(ctx GmvcContext, resp *Response) {
	setDefault(resp)
	setHeader(ctx, resp)
	ctx.HttpResponse().HTML(resp.StatusCode, resp.Body.(string), resp.Model)
}

// StringResponsor string实现
type StringResponsor struct{}

// Response 返回的方法
func (r *StringResponsor) Response(ctx GmvcContext, resp *Response) {
	setDefault(resp)
	setHeader(ctx, resp)

	ctx.HttpResponse().SetHeader("Content-Type", "text/plain")
	ctx.HttpResponse().Status(resp.StatusCode)
	ctx.HttpResponse().Body(strings.NewReader(resp.Body.(string)))
}

var _ Responsor = (*JSONResponsor)(nil)
var _ Responsor = (*HTMLResponsor)(nil)
var _ Responsor = (*StringResponsor)(nil)

func setDefault(resp *Response) {
	if resp.StatusCode == 0 {
		if resp.Body != nil {
			resp.StatusCode = http.StatusOK
		} else {
			resp.StatusCode = http.StatusNoContent
		}
	}
}

func setHeader(ctx GmvcContext, resp *Response) {
	if len(resp.Header) > 0 {
		for k, v := range resp.Header {
			ctx.HttpResponse().SetHeader(k, v)
		}
	}
}
