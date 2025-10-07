## What is gmvc?

Gmvc is a Golang HTTP framework designed to provide a cleanly layered and easy-to-maintain paradigm for web service development.
It offers rich features such as parameter parsing, parameter validation, and customizable response formats, and can be easily integrated with popular web frameworks like Hertz and Gin.

After being refined through dozens of projects and several years of practical use, Gmvc has matured into a stable and capable framework.

## Features

- Cleanly layered architecture
- Easy-to-maintain codebase
- Rich features such as parameter parsing, parameter validation, and customizable response formats
- Easily integrable with popular web frameworks like Hertz and Gin

## Installation

To install Gmvc, simply run:

```
go get github.com/zhengrenjie/gmvc
```

## Getting Started

> [Document: QuickStart](https://github.com/zhengrenjie/gmvc/tree/main/.wiki/1-QuickStart.md)

Let's start with a simplest example:

First, define an gmvc action named `ExampleAction`.


```go

// ExampleAction is a gmvc action.
type ExampleAction struct {
	Ctx context.Context

	Name string `param:"Auto"`
	Age  int    `param:"Auto"`
}

func (a *ExampleAction) Go() (any, error) {
	return struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: a.Name,
		Age:  a.Age,
	}, nil
}
```

And Then, register this action into the gmvc. In this case, we use **Hertz** as the underlying web framework.

```go
func main() {
	// 1. create an gmvc builder
	builder := gmvc_hertz.CreateGmvc4HertzBuilder()

	// 2. create a hertz server
	h := server.Default()

	// 3. wrap the gmvc action to hertz handler
	h.GET("/hello", builder.Wrap(&ExampleAction{}))

	// 4. start the server
	h.Spin()
}
```

Ok, that's all. Now, if you visit `http://127.0.0.1:8888/hello?Name=gmvc&Age=18`, you will get the following response:

```json
{
	"name": "gmvc",
	"age": 18
}
```

## Documentation

For more detailed documentation, please refer to the [wiki](https://github.com/zhengrenjie/gmvc/tree/main/.wiki).

1. [QuickStart](https://github.com/zhengrenjie/gmvc/tree/main/.wiki/1-QuickStart.md)
2. [Parameter Resolver](https://github.com/zhengrenjie/gmvc/tree/main/.wiki/2-Parameter-Resolver.md)
3. [Parameter Binding](https://github.com/zhengrenjie/gmvc/tree/main/.wiki/3-Parameter-Binding.md)
4. [Response Render](https://github.com/zhengrenjie/gmvc/tree/main/.wiki/4-Response-Render.md)



