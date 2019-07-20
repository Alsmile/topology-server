package db

import (
	"topology/db/mongo"
	"topology/db/redis"
)

// Init 初始化数据库连接
func Init() bool {
	return mongo.Init() && redis.Init()
}
