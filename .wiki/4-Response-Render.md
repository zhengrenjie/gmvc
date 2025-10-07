## [WIP] Response Render

Before we go further, let's have a quick look at the response render in gmvc.

In gmvc, there are three ways to make a HTTP response:

1. Return a `any` type value from the `Go` method.
2. Return a [Response](https://github.com/zhengrenjie/gmvc/blob/main/response.go#L26) struct from the `Go` method.
3. Use `Ctx` to call the native HTTP method to make a response.

### Return "Any"
