package tool

import (
	"topology/keys"

	"github.com/kataras/iris/v12"
	"gopkg.in/mgo.v2/bson"
)

// ToolGet 获取用户工具图标列表
func ToolGet(ctx iris.Context) {
	isOperate := ctx.Values().GetBoolDefault("operate", false)
	params := bson.M{"state": 1}
	if isOperate {
		params = bson.M{}
	}

	min := ctx.URLParam("min")
	if min != "" {
		params["base"] = true
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

// ToolsAdd 批量新增
func ToolsAdd(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	data := make([]Tool, 0)
	err := ctx.ReadJSON(&data)
	if err != nil {
		ret["error"] = keys.ErrorParam
		ret["errorDetail"] = err.Error()
		return
	}

	for _, tool := range data {
		Save(&tool, ctx.Values().GetString("uid"), ctx.Values().GetString("username"), ctx.Values().GetBoolDefault("operate", false))
	}
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

// ToolsSave 批量修改
func ToolsSave(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	params := struct {
		IDs   []string
		State int
		Class string
	}{}
	err := ctx.ReadJSON(&params)
	if err != nil {
		ret["error"] = keys.ErrorParam
		ret["errorDetail"] = err.Error()
		return
	}

	objIds := make([]bson.ObjectId, len(params.IDs))
	for i, id := range params.IDs {
		objIds[i] = bson.ObjectIdHex(id)
	}

	data := bson.M{}
	if params.State != 0 {
		data["state"] = params.State
	}
	if params.Class != "" {
		data["class"] = params.Class
	}

	err = Updates(objIds, &bson.M{"$set": data}, ctx.Values().GetString("uid"), ctx.Values().GetString("username"))
	if err != nil {
		ret["error"] = keys.ErrorSave
	}
}

// ToolDel 删除
func ToolDel(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	ids := make([]string, 0)
	err := ctx.ReadJSON(&ids)
	if err != nil {
		ret["error"] = keys.ErrorID
		return
	}

	objIds := make([]bson.ObjectId, len(ids))
	for i, id := range ids {
		objIds[i] = bson.ObjectIdHex(id)
	}

	err = Del(objIds, ctx.Values().GetString("uid"), ctx.Values().GetString("username"))
	if err != nil {
		ret["error"] = keys.ErrorSave
	}
}
