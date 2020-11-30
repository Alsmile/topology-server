package redis

import (
	"strings"
	"time"

	"topology/config"

	"github.com/gomodule/redigo/redis"
	"github.com/mna/redisc"
	"github.com/rs/zerolog/log"
)

// Pool Redis连接池
var Pool *redisc.Cluster

// Init 初始化Redis连接池
func Init() bool {
	Pool = &redisc.Cluster{
		StartupNodes: strings.Split(config.App.Redis.Address, ","),
		DialOptions:  []redis.DialOption{redis.DialConnectTimeout(5 * time.Second)},
		CreatePool:   createPool,
	}

	// 确认连接有效
	// 调用Get()执行RedisPool.Dial连接redis。
	conn := Pool.Get()
	defer conn.Close()

	_, err := conn.Do("SELECT", config.App.Redis.Database)
	if err != nil {
		log.Error().Err(err).Msg("Fail to connect redis.")
	}

	return err == nil
}

func createPool(addr string, opts ...redis.DialOption) (*redis.Pool, error) {
	return &redis.Pool{
		MaxIdle:     5,
		MaxActive:   config.App.Redis.MaxConnections,
		IdleTimeout: time.Duration(config.App.Redis.Timeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr, opts...)

			if err != nil {
				return nil, err
			}
			if config.App.Redis.Password != "" {
				if _, err = c.Do("AUTH", config.App.Redis.Password); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
	}, nil
}
