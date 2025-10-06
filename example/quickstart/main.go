package quickstart

import (
	"context"

	gmvc_hertz "github.com/zhengrenjie/gmvc/adapter/hertz"

	"github.com/cloudwego/hertz/pkg/app/server"
)

type ExampleAction struct {
	Ctx context.Context

	Name string `param"Auto"`
	Age  int    `param"Auto"`
}

func (a *ExampleAction) Go() (any, error) {
	return nil, nil
}

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
