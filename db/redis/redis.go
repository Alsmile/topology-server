package redis

import (
	"time"

	"topology/config"

	"github.com/garyburd/redigo/redis"
	"github.com/rs/zerolog/log"
)

// Pool Redis连接池
var Pool *redis.Pool

// Init 初始化Redis连接池
func Init() bool {
	Pool = &redis.Pool{
		MaxIdle:     config.App.Redis.MaxConnections,
		MaxActive:   config.App.Redis.MaxConnections,
		IdleTimeout: time.Duration(config.App.Redis.Timeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.App.Redis.Address)

			if err != nil {
				return nil, err
			}
			if config.App.Redis.Password != "" {
				if _, err = c.Do("AUTH", config.App.Redis.Password); err != nil {
					c.Close()
					return nil, err
				}
			}

			// Select database.
			_, err = c.Do("SELECT", config.App.Redis.Database)
			return c, nil
		},
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
