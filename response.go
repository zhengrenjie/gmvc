package gmvc

import "net/http"

// RenderType 返回渲染类型
type RenderType int32

const (
	// JSON response, content-type="application/json"
	JSON RenderType = 0

	// HTML 模版渲染, content-type="text/html"
	HTML RenderType = 1

	// String 返回
	String RenderType = 2
)

// Response HTTP返回结果结构体
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
	ctx.JSON(resp.StatusCode, resp.Body)
}

// HTMLResponsor html实现
type HTMLResponsor struct{}

// Response 返回的方法
func (r *HTMLResponsor) Response(ctx GmvcContext, resp *Response) {
	setDefault(resp)
	setHeader(ctx, resp)
	ctx.HTML(resp.Body.(string), resp.StatusCode, resp.Model)
}

// StringResponsor string实现
type StringResponsor struct{}

// Response 返回的方法
func (r *StringResponsor) Response(ctx GmvcContext, resp *Response) {
	setDefault(resp)
	setHeader(ctx, resp)
	ctx.String(resp.StatusCode, resp.Body.(string))
}

func setDefault(resp *Response) {
	if resp.StatusCode == 0 {
		resp.StatusCode = http.StatusOK
	}
}

func setHeader(ctx GmvcContext, resp *Response) {
	if len(resp.Header) > 0 {
		for k, v := range resp.Header {
			ctx.Header(k, v)
		}
	}
}
