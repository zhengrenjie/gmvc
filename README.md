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

Let's start with a simplest example:

First, let's define an gmvc action named `ExampleAction`.


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

And Then, let's register this action into the gmvc. In this case, we use **Hertz** as the underlying web framework.

```go
func main() {
	// 1. create an gmvc builder
	builder := gmvc_hertz.CreateGmvc4HertzBuilder()

	// 2. create a hertz server
	h := server.Default()

	// 3. wrap an gmvc action to hertz handler
	h.GET("/hello", builder.Wrap(&ExampleAction{}))

	// 4. start the server
	h.Spin()
}
```

Ok, that's all. Now, if you visit `/hello?name=gmvc&age=18` (in **Hertz**, the port will be `8888` by default), you will get the following response:

```json
{
	"name": "gmvc",
	"age": 18
}
```

For more detailed documentation, please refer to the [wiki](https://github.com/zhengrenjie/gmvc/wiki).


