package images

import (
	"topology/keys"

	"github.com/kataras/iris/v12"
	"gopkg.in/mgo.v2/bson"
)

// UserImages 获取用户上传的图片列表
// [query] pageIndex - 当前第几页
// [query] pageCount - 每页显示个数
// [query] count - 0，表示不统计总数，返回count为0
// [query] deleted - 查询已删除
// [query] sort - 排序
func UserImages(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	pageIndex, err := ctx.URLParamInt(keys.PageIndex)
	if err != nil {
		ret["error"] = keys.ErrorParamPage
		return
	}
	pageCount, err := ctx.URLParamInt(keys.PageCount)
	if err != nil {
		ret["error"] = keys.ErrorParamPage
		return
	}

	deleted := ctx.URLParam("deleted")

	where := bson.M{
		"userId":    ctx.Values().GetString("uid"),
		"deletedAt": bson.M{"$exists": deleted != ""},
	}

	sort := ctx.URLParam("sort")
	if sort == "" {
		sort = "-createdAt"
	}
	list, count, err := List(&where, sort, pageIndex, pageCount, ctx.URLParam("count") != "0")
	if err != nil {
		ret["error"] = keys.ErrorRead
		ret["errorDetail"] = err.Error()
	}

	ret["list"] = list
	ret["count"] = count
}

// UserImageAdd 新增图片记录
func UserImageAdd(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	data := &Image{}
	err := ctx.ReadJSON(data)
	if err != nil {
		ret["error"] = keys.ErrorParam
		ret["errorDetail"] = err.Error()
		return
	}

	err = Put(data, ctx.Values().GetString("uid"), ctx.Values().GetString("username"))
	if err != nil {
		ret["error"] = keys.ErrorSave
		ret["errorDetail"] = err.Error()
	}

	ret["id"] = data.ID
}

// UserImageDel 删除用户图片记录
func UserImageDel(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	id := ctx.Params().Get("id")
	if !bson.IsObjectIdHex(id) {
		ret["error"] = keys.ErrorID
		return
	}
	err := Del(id, ctx.Values().GetString("uid"))
	if err != nil {
		ret["error"] = keys.ErrorPermission
	}
}
