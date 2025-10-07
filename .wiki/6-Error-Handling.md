## Error Handling

In every web service, global error handling mechanism is always needed. In `golang`, how panic is handled is also a important topic.

You have found that, in previous sections, the error throwing is everywhere, but we don't have discuss how those errors are handled and eventually returned to client.

### Error Handling Mechanism

Gmvc provides a place to register a global error handler. Any error in gmvc's life cycle will be finally caught by this global error handler.

Here is the definition of error handler:

```go
// HandleError is the global error handler.
// If any error occurs during gmvc runtime, it will be catched by this handler.
// This handler convert the error to a response.
type HandleError func(ctx GmvcContext, err error) interface{}
```

This method has two input parameters:

- `ctx` is the gmvc context.
- `err` is the error that occurs during gmvc runtime.

And one output parameter:

- `interface{}` is the response that will be returned to client.

Look here, `HandleError` finally returns a `any` type value to gmvc framework, and this value will be handled exactly same as the value returned by `Go`! Which means, any feature supported by [Response Render](https://github.com/zhengrenjie/gmvc/tree/main/.wiki/4-Response-Render.md), is also supported to `HandleError`.

So here is the magic: every error in the runtime will be eventually converted to a standard gmvc response and returned to client by gmvc `Response Render`.

### Error Handler Registration

To register a global error handler, you need to call `RegisterErrorHandler` method in `Gmvc` instance.

Example:

```go
var MyErrorHandler = func(ctx GmvcContext, err error) interface{} {
	// ... custom error handling logic
    if errors.Is(err, error1) {
        // do something
    } else if errors.Is(err, error2) {
        // do something else
    }

	return nil
}

// main
func main() {
	builder := gmvc_hertz.CreateGmvc4HertzBuilder()

    // set global error handler to gmvc.
	builder.SetErrorHandler(MyErrorHandler)

	h := server.Default()
	h.GET("/hello", builder.Wrap(&ExampleAction{}))
	h.Spin()
}
```

## What's next?

- [Middleware](https://github.com/zhengrenjie/gmvc/tree/main/.wiki/7-Middleware.md)
