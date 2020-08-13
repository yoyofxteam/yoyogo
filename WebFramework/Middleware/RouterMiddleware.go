package Middleware

import (
	"github.com/yoyofx/yoyogo/WebFramework/Context"
	"github.com/yoyofx/yoyogo/WebFramework/Router"
	"strings"
)

//var ReqFuncMap = make(map[string]func(ctx * YoyoGo.HttpContext))

type RouterMiddleware struct {
	RouterBuilder Router.IRouterBuilder
}

func NewRouter(builder Router.IRouterBuilder) *RouterMiddleware {
	return &RouterMiddleware{RouterBuilder: builder}
}

func (router *RouterMiddleware) Inovke(ctx *Context.HttpContext, next func(ctx *Context.HttpContext)) {
	var handler func(ctx *Context.HttpContext)
	handler = router.RouterBuilder.Search(ctx, strings.Split(ctx.Input.Request.URL.Path, "/")[1:], ctx.Input.RouterData)
	if handler != nil {
		handler(ctx)
	}
	next(ctx)

}
