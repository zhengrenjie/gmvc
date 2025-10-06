package gmvc_hertz

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/zhengrenjie/gmvc"
)

type Gmvc4HertzBuilder struct {
	*gmvc.GmvcBuilder
}

func (g *Gmvc4HertzBuilder) Wrap(action gmvc.Action, mdw ...gmvc.IMiddleware) app.HandlerFunc {
	return Wrap(g.BuildAction(action, mdw...))
}

func CreateGmvc4HertzBuilder() *Gmvc4HertzBuilder {
	builder := gmvc.CreateGmvcBuilder(gmvc.DefineAuto(gmvc.QuerySrc, gmvc.FormSrc))
	return &Gmvc4HertzBuilder{
		GmvcBuilder: builder,
	}
}

func Wrap(handler gmvc.HandlerFunc) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		adapter := AcquireHertzContext(c, ctx)
		handler(adapter)
	}
}
