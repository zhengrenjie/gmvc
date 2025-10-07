## Parameter Resolver

### Primary Type Resolver

Gmvc has built-in primary type resolver.

```go
type ExampleAction struct {
	Ctx context.Context

	String  string  `param:"Auto"`
	Int64   int64   `param:"Auto"`
	PInt32  *int32  `param:"Auto"`
	Float64 float64 `param:"Auto"`

	Bool bool `param:"Auto"`

	IntSlice     []int     `param:"Auto"`
	StringSlice  []string  `param:"Auto"`
	StringPSlice []*string `param:"Auto"`
}

func (a *ExampleAction) Go() (any, error) {
	// ... ignore response here
    return nil, nil
}
```

The supported primary types are:

| Supported Primary Types      | Support Pointer | Support Array | Support Array of Pointer |
| ----------- | ----------- | ----------- | ----------- |
| `string`    | ✅ | ✅ | ✅ |
| `int`       | ✅ | ✅ | ✅ |
| `int8`      | ✅ | ✅ | ✅ |
| `int16`     | ✅ | ✅ | ✅ |
| `int32`     | ✅ | ✅ | ✅ |
| `int64`     | ✅ | ✅ | ✅ |
| `uint`      | ✅ | ✅ | ✅ |
| `uint8`     | ✅ | ✅ | ✅ |
| `uint16`    | ✅ | ✅ | ✅ |
| `uint32`    | ✅ | ✅ | ✅ |
| `uint64`    | ✅ | ✅ | ✅ |
| `float64`   | ✅ | ✅ | ✅ |
| `bool`      | ✅ | ✅ | ✅ |

For `bool` type, the resolver will convert the following values to `true`:

- `true`
- `TRUE`
- `True`

And the following values to `false`:

- `false`
- `FALSE`
- `False`

For `Array` type, the built-in resolver will split the value by comma `,` and then convert each item to the corresponding type.

### Built-in Struct Resolver

#### Json Resolver

Gmvc has one built-in struct resolver named `Json`.

```go
type User struct {
	Name string
	Age  int
}

type ExampleAction struct {
	Ctx context.Context

	User User `param:"Auto" resolver:"Json"`
}

func (a *ExampleAction) Go() (any, error) {
	// ... ignore response here
    return nil, nil
}
```

When you have a struct to unmarshal, you can use the `Json` resolver to unmarshal the JSON string. That would be useful when you want to parse a JSON parameter.

Notice, the `Json` resolver does not mean that it will take `application/json` as the content type and get data from HTTP-Body. Instead, it does not care about the content type. It will just try to unmarshal the origin value got from HTTP parameter to the JSON string.

You might ask where does gmvc get the origin value from. The answer is it is determined by the `param` tag. In the above example, we use `param:"Auto"` to let gmvc get the value from multiple place of HTTP text.

Also, in this case, the `User` parameter also takes `*User`. For `Array`, as we use `Json` resolver, it will just use `json.Unmarshal` to resolver the array type, instead of spliting the value by comma `,`, which only works for the primary-type resolver.

#### Context Resolver

You may have noticed that in every (almost) examples, there is a parameter named `Ctx`. It is always recommended to add it to your Action. 

You can have three types of context:

- `context.Context`
- `gmvc.GmvcContext`
- `hertz.RequestContext` or other underlying web framework's context.

If you want to use the origin HTTP context, it is recommended to use `gmvc.GmvcContext`, where you can get ride of the underlying web framework's context and make your business logic more clear and decoupled with any certain underlying web framework.

If you only want to do some processing control, like timeout controlling, or logging & tracing , you just need `context.Context`.

You should use `hertz.RequestContext` or any certain underlying web framework's context only when you really need it. Or you are coupling your business logic with a certain underlying web framework.


### Custom Resolver

Also you can define your own resolver, in case you want to resolve the value in a special way. You can do like this:

```go

// define custom resolver
var MyResolver = func(ctx gmvc.GmvcContext, fieldMeta *gmvc.ParamMeta, origin string) (interface{}, error) {
	user := User{}

	err := json.Unmarshal([]byte(origin), &user)
	if err != nil {
		return nil, err
	}

    // use "虚岁" to represent the age
	user.Age += 1
	return user, nil
}

type User struct {
	Name string
	Age  int
}

type ExampleAction struct {
	Ctx context.Context

	User User `param:"Auto" resolver:"MyResolver"` // indicate use MyResolver to resolve User parameter
}

func (a *ExampleAction) Go() (any, error) {
	// ... ignore response here
    return nil, nil
}

// main
func main() {
	builder := gmvc_hertz.CreateGmvc4HertzBuilder()

    // register your resolver to gmvc.
	builder.RegisterResolver("MyResolver", MyResolver)

	h := server.Default()

	// wrap an gmvc action to hertz handler
	h.GET("/hello", builder.Wrap(&ExampleAction{}))
	h.Spin()
}

```

In this case, you define a custom resolver named `MyResolver`. It will unmarshal the JSON string to `User` struct and then add `1` to the `Age`.

#### Resolver Definition

The `Resolver` function is define like:

```go
type Resolver func(ctx gmvc.GmvcContext, fieldMeta *gmvc.ParamMeta, origin string) (interface{}, error)
```

The three parameters of `Resolver` function are:

- `ctx`: the gmvc context, you can get the origin HTTP context from it.
- `fieldMeta`: the metadata of the Action field, you can get the reflecting info from it.
- `origin`: the origin value got from HTTP parameter.

The return value of `Resolver` function is:

- `interface{}`: the resolved value, **it must be exactly the same type as the field** (eg.if the field in Action is pointer type, the resolved value must be pointer type as well).
- `error`: the error occurred during the resolution.

#### Resolver Registration

After you define your resolver, you need to register it to gmvc. You can do like this:

```go
builder.RegisterResolver("MyResolver", MyResolver)
```

The first parameter is the resolver name, which you can use in the `resolver` tag. The second parameter is the resolver function.

## What's next?

- [Parameter Binding](https://github.com/zhengrenjie/gmvc/tree/main/.wiki/3-Parameter-Binding.md)
