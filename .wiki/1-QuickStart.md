# QuickStart

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


What happens here?

First, we define an action `ExampleAction` with two parameters `Name` and `Age` (actually it's three, but we don't discuss `Ctx` here), and a method `Go`. In gmvc, you need to use `param:""` tag to indicate which parameter should be bind from the HTTP request. The `Auto` means gmvc will bind the parameter automatically without indicating **from where** (e.g. query, path, form, etc.). That will be useful when you need to support one parameter in multiple ways.

And then, in `ExampleAction` you should provider `Go` method, with the method signature `func Go() (any, error)`. This method make sure the `ExampleAction` has implemented the [gmvc.Action](https://github.com/zhengrenjie/gmvc/blob/main/interface.go#L143) interface. Gmvc will test it, when you try to register this action into the gmvc it will check if the action implements the interface, to make sure it wouldn't startup if there are any violations. In the runtime, when the HTTP request comes in, gmvc will create an instance of `ExampleAction` and bind the parameters from the request via gmvc tags, then call `Go` method to get the response.

You already have the Action, now let's register it into the web framework. In this case, we use **Hertz** as the underlying web framework.

Gmvc is not an web framework, which means it wouldn't start the server by itself, or have any underlying connections, or implement any web protocols. Instead, it provides a way to wrap the gmvc action into the web framework handler. In this case, we use `builder.Wrap` to wrap the `ExampleAction` into a hertz handler.

As you see, gmvc has already prepare the "adapter" for hertz, you can get a instance of `GmvcBuilder` for Hertz by calling `gmvc_hertz.CreateGmvc4HertzBuilder()`. And then, you can use `builder.Wrap` to wrap the `ExampleAction` into a hertz handler.

`GmvcBuilder` is the main interface for everyone who are using gmvc, it provides a lot of methods to help you to register the actions, middlewares, checkers, resolvers, and other configurations. We will discuss them in detail in the following sections.

Now, everything is ready, you have the knowledge of how to develop a gmvc program and have the basic understanding of how gmvc works, and you can start this simplest example or play with gmvc in your own project.















