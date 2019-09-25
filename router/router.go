package router

import (
	"strconv"
	"strings"

	"topology/config"
	"topology/images"
	"topology/middlewares"
	"topology/topology"
	"topology/websocket"

	"github.com/kataras/iris"
)

// Listen 监听路由
func Listen() {
	route := iris.New()
	route.OnErrorCode(iris.StatusNotFound, NotFound)

	route.Use(middlewares.Usr)

	// 拓扑图模块
	topology.Route(route)

	// 用户图库模块
	images.Route(route)

	websocket.Route(route)

	route.StaticWeb("/", "../web")

	// 监听
	port := strconv.Itoa(int(config.App.Port))
	route.Run(
		iris.Addr(":"+port),
		// skip err server closed when CTRL/CMD+C pressed:
		iris.WithoutServerError(iris.ErrServerClosed),
		// enables faster json serialization and more:
		iris.WithOptimizations,
	)
}

// Index 首页静态文件
func Index(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
	ctx.ServeFile("../web/index.html", false)
}

// NotFound 404
func NotFound(ctx iris.Context) {
	if strings.HasPrefix(ctx.Path(), "/api/") {
		ret := make(map[string]interface{})
		ret["error"] = "请求错误（Not found）：" + ctx.Path()
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(ret)
	} else {
		ctx.StatusCode(iris.StatusFound)
		Index(ctx)
	}
}
