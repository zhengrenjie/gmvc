## Parameter Binding

I believe you have already got some idea from the previous section of how parameter binding works in gmvc. Actually, the parameter binding happens before the `Resolver` in gmvc. In this section, we will discuss the details of parameter binding.

First, a parameter binding rule is defined as follows:

- If the field must have a `param` tag, otherwise the value will be ignored.

### param tag

The `param` tag formats as: `param:"<Location>,<Name>,<RecursiveMark>"`.

- `<Location>`: [Required] the location of the parameter. It can be `Query`, `Form`, `Body`, `Header`, `Path`, `Ctx`, or `Auto`.
- `<Name>`: [Optional] to rename of the parameter. If it is empty, the field name will be used.
- `<RecursiveMark>`: [Optional] the recursive mark. It can be `Recursive` or empty. If it is `Recursive`, the parameter will be parsed recursively.

The three parts does not have any order.

### Location Tag

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

In most cases, `Auto` is enough. But when you want to bind a parameter in a specific location, or when you start to care about the order of how `Auto` lookups the parameter, you should not use `Auto`.

### Rename Tag

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

### Recursive Tag

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

### Default Value Tag

You can set a default value for a parameter by using `default` tag.

Example:

```go
// ExampleAction is a gmvc action.
type ExampleAction struct {
	Ctx context.Context

    PageSize int `param:"Query" default:"10"`
}
```

Notice, the default value will only be set when there is not such Key `PageSize`, in this case, in HTTP Query. If the key exist even with empty value, like "/?PageSize=", the default value will not be set.

**Strictly speaking**, `default` tag is an `Resolver` function, as it is to set the value into the Action, instead of get value from HTTP parameter. But I still find it is more clear to discuss it here.

You might have got some ideas of the **"Life Cycle"** in gmvc right? The parameter binding happens before the `Resolver`. We will disscuss this later.

## What's next?

- [Response Render](https://github.com/zhengrenjie/gmvc/tree/main/.wiki/4-Response-Render.md)
