package router

import (
	"strconv"
	"strings"

	"topology/cms"
	"topology/config"
	"topology/images"
	"topology/middlewares"
	"topology/tool"
	"topology/topology"
	"topology/websocket"

	"github.com/kataras/iris/v12"
)

// Listen 监听路由
func Listen() {
	route := iris.New()
	route.OnErrorCode(iris.StatusNotFound, NotFound)

	route.Use(middlewares.Usr)

	topology.Route(route)
	images.Route(route)
	cms.Route(route)
	tool.Route(route)

	websocket.Route(route)

	route.HandleDir("/", "../web")

	// p := pprof.New()
	// route.Any("/debug/pprof", p)
	// route.Any("/debug/pprof/{action:path}", p)

	// 监听
	route.Listen(":" + strconv.Itoa(int(config.App.Port)))
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
