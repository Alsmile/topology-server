package topology

import (
	"topology/middlewares"

	"github.com/kataras/iris/v12"
)

// Route file模块路由
func Route(route *iris.Application) {
	route.Get("/api/topology/:id", GetTopology)
	route.Get("/api/topologies", Topologies)

	routeUser := route.Party("/api/user", middlewares.Auth)
	routeUser.Get("/topologies", UserTopologies)
	routeUser.Post("/topology", UserTopologyAdd)
	routeUser.Put("/topology", UserTopologySave)
	routeUser.Patch("/topology", UserTopologyPatch)
	routeUser.Delete("/topology/:id", UserTopologyDel)
	routeUser.Post("/topology/restore/:id", UserTopologyRestore)

	routeUser.Get("/topology/histories", TopologyHistories)
	routeUser.Patch("/topology/history", TopologyHistoryPatch)
	routeUser.Delete("/topology/history/:id", TopologyHistoryDel)

	routeUser.Get("/stars", UserStars)
	routeUser.Post("/star", UserStarAdd)
	routeUser.Delete("/star/:id", UserStarDel)
	routeUser.Post("/star/ids", UserStarIDs)

	routeUser.Get("/statistics", UserStatistics)
}
