package main

import (
	"context"
	"encoding/json"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/zhengrenjie/gmvc"
	gmvc_hertz "github.com/zhengrenjie/gmvc/adapter/hertz"
)

var MyResolver = func(ctx gmvc.GmvcContext, fieldMeta *gmvc.ParamMeta, origin string) (interface{}, error) {
	user := User{}

	err := json.Unmarshal([]byte(origin), &user)
	if err != nil {
		return nil, err
	}

	user.Age += 1
	return user, nil
}

type User struct {
	Name string
	Age  int
}

type ExampleAction struct {
	Ctx context.Context

	String  string  `param:"Auto"`
	Int     int     `param:"Auto"`
	Int8    int8    `param:"Auto"`
	Int16   int16   `param:"Auto"`
	Int32   int32   `param:"Auto"`
	Int64   int64   `param:"Auto"`
	PInt32  *int32  `param:"Auto"`
	Float64 float64 `param:"Auto"`

	Bool bool `param:"Auto"`

	IntSlice     []int     `param:"Auto"`
	StringSlice  []string  `param:"Auto"`
	StringPSlice []*string `param:"Auto"`

	User   User `param:"Auto" resolver:"Json"`
	MyUser User `param:"Auto" resolver:"MyResolver"`
}

func (a *ExampleAction) Go() (any, error) {
	return struct {
		String  string  `json:"string"`
		Int32   int32   `json:"int32"`
		PInt32  *int32  `json:"pint32"`
		Float64 float64 `json:"float64"`

		StringSlice  []string  `json:"string_slice"`
		StringPSlice []*string `json:"string_p_slice"`

		Bool bool `json:"bool"`

		User   User `json:"user"`
		MyUser User `json:"my_user"`
	}{
		String:       a.String,
		Int32:        a.Int32,
		PInt32:       a.PInt32,
		Float64:      a.Float64,
		StringSlice:  a.StringSlice,
		StringPSlice: a.StringPSlice,
		Bool:         a.Bool,
		User:         a.User,
		MyUser:       a.MyUser,
	}, nil
}

func main() {
	builder := gmvc_hertz.CreateGmvc4HertzBuilder()
	builder.RegisterResolver("MyResolver", MyResolver)

	h := server.Default()

	// wrap an gmvc action to hertz handler
	h.GET("/hello", builder.Wrap(&ExampleAction{}))
	h.Spin()
}
