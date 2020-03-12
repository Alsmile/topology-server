package session

import (
	"net/http"

	"github.com/kataras/iris/v12"
	"topology/db/redis"
	"topology/utils"
)

const (
	// Name session name
	Name = "sid"
	// MaxAge 10分钟，只用于短期存储
	MaxAge = 10 * 60
)

// GetsessionID 获取一个sessionID，没用时，自动生成
func GetsessionID(ctx iris.Context) (sessionID string) {
	sessionID = ctx.GetCookie(Name)
	if sessionID == "" {
		sessionID = utils.GetGUID()
		cookie := &http.Cookie{}
		cookie.Path = "/"
		cookie.HttpOnly = true
		cookie.Name = Name
		cookie.Value = sessionID
		ctx.SetCookie(cookie)
	}

	return
}

// SetSession 设置name=val的session
func SetSession(ctx iris.Context, name string, val interface{}) error {
	sessionID := GetsessionID(ctx)

	if sessionID == "" {
		sessionID = GetsessionID(ctx)
	}

	redisConn := redis.Pool.Get()
	defer redisConn.Close()
	_, err := redisConn.Do("SETEX", sessionID+"."+name, MaxAge, val)
	return err
}

// GetSession 通过name获取session值
func GetSession(ctx iris.Context, name string) (val interface{}, err error) {
	sessionID := GetsessionID(ctx)
	if sessionID == "" {
		return
	}

	redisConn := redis.Pool.Get()
	defer redisConn.Close()

	val, err = redisConn.Do("GET", sessionID+"."+name)
	return
}
