## Parameter Binding

I believe you have already got some idea of parameter binding from the previous section. Actually, the parameter binding happens before the `Resolver` in gmvc. In this section, we will discuss the details of parameter binding.

First, a parameter binding rule is defined as follows:

- If the field must have a `param` tag, otherwise the value will be ignored.

### param tag

The `param` tag formats as: `param:"<Location>,<Name>,<RecursiveMark>"`.

- `<Location>`: [Required] the location of the parameter. It can be `Query`, `Form`, `Body`, `Header`, `Path`, `Ctx`, or `Auto`.
- `<Name>`: [Optional] to rename of the parameter. If it is empty, the field name will be used.
- `<RecursiveMark>`: [Optional] the recursive mark. It can be `Recursive` or empty. If it is `Recursive`, the parameter will be parsed recursively.

The three parts does not have any order.

These are some examples:

- `param:"Query"`
- `param:"Query,X-Real-Ip"`
- `param:"Query,Name,Recursive"`
- `param:"Form,Name"`
- `param:"Body,Name"`
- `param:"Header,Name"`
- `param:"Path,Name"`
- `param:"Ctx,Name"`
- `param:"Auto,Name"`

In most cases, `Auto` is enough. But when you want to bind a parameter in a specific location, or when you start to care about the order of how `Auto` lookup the parameter, you should not use `Auto`.

### Rename

Sometimes you must need to rename the parameter, eg. when you want to get Header from HTTP (because in most cases, the HTTP Header key is not a valid Go identifier).

Example:

```go
// ExampleAction is a gmvc action.
type ExampleAction struct {
	Ctx context.Context

    // X-Real-Ip is not a valid Go identifier, so we rename it to RealIp.
	RealIp string `param:"Header,X-Real-Ip"`
}
```

### Recursive

Sometimes you may want to group some parameters together, or you have some common parameter pairs for all over the Actions. In those cases, you may need `Recursive` mark.

Example:

```go

// The common page parameter pair.
type Page struct {
    PageNumber int `param:"Query"`
    PageSize int `param:"Query" default:"10"`
}

// ExampleAction is a gmvc action.
type ExampleAction struct {
	Ctx context.Context

    Page Page `param:"Recursive"`
}
```