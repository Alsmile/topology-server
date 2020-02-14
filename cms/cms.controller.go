package cms

import (
	"topology/keys"

	"github.com/kataras/iris"
	"gopkg.in/mgo.v2/bson"
)

// CmsGet 获取指定配置内容
func CmsGet(ctx iris.Context) {
	id := ctx.Params().Get("id")
	if !bson.IsObjectIdHex(id) {
		ctx.JSON(bson.M{
			"error": keys.ErrorID,
		})
		return
	}

	data, err := GetCmsByID(id)
	if err != nil {
		ctx.JSON(bson.M{
			"error":       keys.ErrorRead,
			"errorDetail": err.Error(),
		})
		return
	}
	ctx.JSON(data)
}

// CmsAdd 新增配置内容
func CmsAdd(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	data := &Cms{}
	err := ctx.ReadJSON(data)
	if err != nil {
		ret["error"] = keys.ErrorParam
		ret["errorDetail"] = err.Error()
		return
	}

	data.ID = ""
	err = Save(data, ctx.Values().GetString("uid"), ctx.Values().GetString("username"))
	if err != nil {
		ret["error"] = keys.ErrorSave
		ret["errorDetail"] = err.Error()
	}

	ret["id"] = data.ID
}

// CmsSave 保存修改
func CmsSave(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	data := &Cms{}
	err := ctx.ReadJSON(data)
	if err != nil || data.ID == "" {
		ret["error"] = keys.ErrorParam
		if err != nil {
			ret["errorDetail"] = err.Error()
		}
		return
	}

	err = Save(data, ctx.Values().GetString("uid"), ctx.Values().GetString("username"))
	if err != nil {
		ret["error"] = keys.ErrorSave
		ret["errorDetail"] = err.Error()
	}
	ret["id"] = data.ID
}

// CmsDel 删除
func CmsDel(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	id := ctx.Params().Get("id")
	if !bson.IsObjectIdHex(id) {
		ret["error"] = keys.ErrorID
		return
	}
	err := Del(id, ctx.Values().GetString("uid"), ctx.Values().GetString("username"))
	if err != nil {
		ret["error"] = keys.ErrorPermission
	}
}
