package images

import (
	"topology/middlewares"

	"github.com/kataras/iris/v12"
)

// Route file模块路由
func Route(route *iris.Application) {
	routeUser := route.Party("/api/user", middlewares.Auth)
	routeUser.Get("/images", UserImages)
	routeUser.Post("/image", UserImageAdd)
	routeUser.Delete("/image/:id", UserImageDel)
}
