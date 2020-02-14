package middlewares

import (
	"fmt"
	"strings"
	"time"

	"topology/config"
	"topology/keys"
	"topology/utils"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"github.com/rs/zerolog/log"
)

// ParseJwt 解析JWT
func ParseJwt(ctx iris.Context, data string) error {
	// jwt校验
	token, err := jwt.Parse(data, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("签名方法错误: %v", token.Header["alg"])
		}
		return []byte(config.App.Jwt), nil
	})

	if err != nil {
		log.Error().
			Err(err).
			Str("func", "middlewares.Usr").
			Str("token", data).
			Str("jwt", config.App.Jwt).
			Str("remoteAddr", ctx.RemoteAddr()).
			Msg("Jwt parse error.")
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Warn().
			Str("func", "middlewares.Usr").
			Str("token", data).
			Str("jwt", config.App.Jwt).
			Str("remoteAddr", ctx.RemoteAddr()).
			Msg("Jwt invalid.")
		return err
	}

	// 设置uid和role
	uid := utils.String(claims["uid"])
	if uid != "" {
		ctx.Values().Set("uid", uid)
		ctx.Values().Set("username", utils.String(claims["username"]))
		ctx.Values().Set("role", utils.String(claims["role"]))
		ctx.Values().Set("vip", utils.Int(claims["vip"]))
		ctx.Values().Set("vipExpiry", utils.Int64(claims["vipExpiry"]))
	}

	return nil
}

// Usr 解析用户身份
func Usr(ctx iris.Context) {
	// 获取header
	data := ctx.GetHeader("Authorization")
	if data == "" {
		ctx.Next()
		return
	}

	ParseJwt(ctx, data)

	ctx.Next()
}

// Auth 身份认证中间件
func Auth(ctx iris.Context) {
	if ctx.Values().GetString("uid") != "" {
		ctx.Next()
		return
	}

	ctx.StatusCode(iris.StatusUnauthorized)
	ret := make(map[string]interface{})
	ret["error"] = keys.ErrorNeedSign
	ctx.JSON(ret)
}

// Vip 获取vip身份
func Vip(ctx iris.Context) uint8 {
	vip, _ := ctx.Values().GetUint8("vip")
	if vip < 1 {
		return vip
	}

	vipExpiry, _ := ctx.Values().GetInt64("vipExpiry")

	if vipExpiry-time.Now().Unix() < -86400 {
		vip = 0
	}

	return vip
}

// Operater 必须是运营人员
func Operater(ctx iris.Context) {
	if !strings.Contains(ctx.Values().GetString("role"), "operation") {
		log.Warn().Str("Illegal access", ctx.Values().GetString("uid")).Msg("auth.Operater")

		return
	}

	ctx.Next()
}
