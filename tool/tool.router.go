package tool

import (
	"topology/middlewares"

	"github.com/kataras/iris"
)

// Route file模块路由
func Route(route *iris.Application) {
	route.Get("/api/tools", ToolGet)

	routeUser := route.Party("/api/user", middlewares.Auth)
	routeUser.Post("/tool", ToolAdd)
	routeUser.Put("/tool", ToolSave)
	routeUser.Delete("/tool/:id", ToolDel)

	routeOperate := route.Party("/api/operate", middlewares.Auth, middlewares.Operater, func(ctx iris.Context) {
		ctx.Values().Set("operate", true)
		ctx.Next()
	})
	routeOperate.Get("/tools", ToolGet)
	routeOperate.Post("/tool", ToolAdd)
	routeOperate.Put("/tool", ToolSave)
	routeOperate.Delete("/tool/:id", ToolDel)
}
