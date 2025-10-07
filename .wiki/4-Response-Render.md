## [WIP] Response Render

Before we go further, let's have a quick look at the response render in gmvc.

In gmvc, there are three ways to make a HTTP response:

1. Return a `any` type value from the `Go` method.
2. Return a [Response](https://github.com/zhengrenjie/gmvc/blob/main/response.go#L26) struct from the `Go` method.
3. Use `Ctx` to call the native HTTP method to make a response.

### Return "Any"

If you return a `any` type value from the `Go` method, gmvc will render it to JSON and set the `Content-Type` header to `application/json` with status code `200 OK`.

This might be the most common way to make a response.

### Return "Response"

```go
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
```

[Response](https://github.com/zhengrenjie/gmvc/blob/main/response.go#L26) is a build-in struct for several common usage:

- JSON: make a `application/json` response
- HTML: make a `text/html` response
- String: make a `text/plain` response

By indicating the `RenderType`, gmvc will use the corresponding responsor to return the response.

In each way, gmvc will use the `StatusCode` as the response status code, and will merge the `Header` with the default headers.

#### JSON Response

Json response might be the most common way to make a response in a web-server. In this case, such response will be rendered:

1. Use `json.Marshal` to marshal `Response.Body` to JSON.
2. Set `Content-Type` header to `application/json`.


#### HTML Response

HTML response will be rendered as follows:

1. `Response.Body` must be a `string` type which indicates the HTML template.
2. Use `html/template` to render `Response.Body` with `Response.Model`.
3. Set `Content-Type` header to `text/html`.

#### String Response

String response will be rendered as follows:

1. `Response.Body` must be a `string` type.
2. Set `Content-Type` header to `text/plain`.

#### Custom Response

Also, gmvc supports to extend the `RenderType` by implementing the `Responsor` interface, or override the default (JSON/HTML/String) responsor.

Example:

```go
var _ Responsor = (*CustomResponsor)(nil)

const MyRender RenderType = 2

type CustomResponsor struct{}

func (r *CustomResponsor) Response(ctx GmvcContext, resp *Response) {
	// ... custom render logic
	return nil
}

// main
func main() {
	builder := gmvc_hertz.CreateGmvc4HertzBuilder()

    // register your render to gmvc.
	builder.RegisterResponsor(MyRender, &CustomResponsor{})

	h := server.Default()

	// wrap an gmvc action to hertz handler
	h.GET("/hello", builder.Wrap(&ExampleAction{}))
	h.Spin()
}
```

### Return by HTTP Context

This is the most native way to write response in HTTP, but also, you should be very careful about what you are doing. In the life cycle of gmvc, you will see in many different stages that you can do response, even in different way. So, again, make sure you know how it works.

But still, this way will be the most flexible way to make a response.

Example:

```go
type ExampleAction struct {
	Ctx gmvc.GmvcContext
}

func (a *ExampleAction) Go() (any, error) {
    // use HTTP Context to write response
    a.Ctx.HttpResponse().Status(http.StatusOK)
    a.Ctx.HttpResponse().Body(strings.NewReader("Hello, World!"))
    a.Ctx.HttpResponse().SetHeader("Content-Type", "text/plain")


	// return no thing here, or it will be conflict with HTTP Context response.
    return nil, nil
}
```

## What's next?

- [Parameter Checker](https://github.com/zhengrenjie/gmvc/tree/main/.wiki/5-Parameter-Checker.md)
