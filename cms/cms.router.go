package cms

import (
	"topology/middlewares"

	"github.com/kataras/iris"
)

// Route file模块路由
func Route(route *iris.Application) {
	route.Get("/api/cms", CmsGet)

	routeUser := route.Party("/api/operate", middlewares.Auth, middlewares.Operater)
	routeUser.Post("/cms", CmsAdd)
	routeUser.Put("/cms", CmsSave)
	routeUser.Delete("/cms/:id", CmsDel)
}
