package tool

import (
	"topology/keys"

	"github.com/kataras/iris/v12"
	"gopkg.in/mgo.v2/bson"
)

// ToolGet 获取用户工具图标列表
func ToolGet(ctx iris.Context) {
	isOperate := ctx.Values().GetBoolDefault("operate", false)
	params := bson.M{}
	if isOperate {
		params["shared"] = true
	} else {
		params["$or"] = []bson.M{
			bson.M{"userId": ctx.Values().GetString("uid")},
			bson.M{"shared": true},
		}
	}

	data, _, err := List(&params, 0, 0, false)
	if err != nil {
		ctx.JSON(bson.M{
			"error":       keys.ErrorRead,
			"errorDetail": err.Error(),
		})
		return
	}
	ctx.JSON(data)
}

// ToolAdd 新增
func ToolAdd(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	data := &Tool{}
	err := ctx.ReadJSON(data)
	if err != nil {
		ret["error"] = keys.ErrorParam
		ret["errorDetail"] = err.Error()
		return
	}

	data.ID = ""
	err = Save(data, ctx.Values().GetString("uid"), ctx.Values().GetString("username"), ctx.Values().GetBoolDefault("operate", false))
	if err != nil {
		ret["error"] = keys.ErrorSave
		ret["errorDetail"] = err.Error()
	}

	ret["id"] = data.ID
}

// ToolSave 修改
func ToolSave(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	data := &Tool{}
	err := ctx.ReadJSON(data)
	if err != nil || data.ID == "" {
		ret["error"] = keys.ErrorParam
		if err != nil {
			ret["errorDetail"] = err.Error()
		}
		return
	}

	err = Save(data, ctx.Values().GetString("uid"), ctx.Values().GetString("username"), ctx.Values().GetBoolDefault("operate", false))
	if err != nil {
		ret["error"] = keys.ErrorSave
		ret["errorDetail"] = err.Error()
	}
	ret["id"] = data.ID
}

// ToolDel 删除
func ToolDel(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	id := ctx.Params().Get("id")
	if !bson.IsObjectIdHex(id) {
		ret["error"] = keys.ErrorID
		return
	}
	err := Del(id, ctx.Values().GetString("uid"), ctx.Values().GetString("username"), ctx.Values().GetBoolDefault("operate", false))
	if err != nil {
		ret["error"] = keys.ErrorPermission
	}
}
