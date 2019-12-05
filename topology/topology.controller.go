package topology

import (
	"time"
	"topology/keys"
	"topology/middlewares"
	"topology/utils"

	"github.com/kataras/iris"
	"gopkg.in/mgo.v2/bson"
)

// GetTopology 获取指定拓扑图
func GetTopology(ctx iris.Context) {
	id := ctx.Params().Get("id")
	if !bson.IsObjectIdHex(id) {
		ctx.JSON(bson.M{
			"error": keys.ErrorID,
		})
		return
	}

	fileID := ctx.URLParam("fileId")
	if fileID != "" && !bson.IsObjectIdHex(fileID) {
		ctx.JSON(bson.M{
			"error": keys.ErrorID,
		})
		return
	}

	if fileID == "" {
		data, err := GetTopologyByID(id, ctx.Values().GetString("uid"))
		if err != nil {
			ctx.JSON(bson.M{
				"error":       keys.ErrorRead,
				"errorDetail": err.Error(),
			})
			return
		}
		ctx.JSON(data)
		return
	}

	data, err := GetHistory(id, fileID, ctx.Values().GetString("uid"))
	if err != nil {
		ctx.JSON(bson.M{
			"error":       keys.ErrorRead,
			"errorDetail": err.Error(),
		})
		return
	}
	ctx.JSON(data)
}

// Topologies 获取已分享的拓扑图
// [query] pageIndex - 当前第几页
// [query] pageCount - 每页显示个数
// [query] count - 0，表示不统计总数，返回count为0
// [query] name - 搜索name
// [query] desc - 搜索desc
// [query] text - 搜索name和desc
// [query] user - 搜索username
// [query] createdStart
// [query] createdEnd
// [query] updatedStart
// [query] updatedEnd
// [query] sort - 排序
func Topologies(ctx iris.Context) {
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

	where := bson.M{
		"shared":    true,
		"deletedAt": bson.M{"$exists": false},
	}
	name := ctx.URLParam("name")
	if name != "" {
		where["name"] = bson.M{"$regex": name, "$options": "$i"}
	}
	desc := ctx.URLParam("desc")
	if desc != "" {
		where["desc"] = bson.M{"$regex": desc, "$options": "$i"}
	}
	text := ctx.URLParam("text")
	if text != "" {
		where["$or"] = []bson.M{
			bson.M{"$regex": name, "$options": "$i"},
			bson.M{"$regex": desc, "$options": "$i"},
		}
	}

	user := ctx.URLParam("user")
	if user != "" {
		where["username"] = bson.M{"$regex": user, "$options": "$i"}
	}

	createdTime := bson.M{}
	createdStart, _ := ctx.URLParamInt64("createdStart")
	if createdStart > 0 {
		createdTime["$gte"] = time.Unix(createdStart, 0)
	}
	createdEnd, _ := ctx.URLParamInt64("createdEnd")
	if createdEnd > 0 {
		createdTime["$lte"] = time.Unix(createdEnd, 0)
	}
	if createdTime["$gte"] != nil || createdTime["$lte"] != nil {
		where["createdAt"] = createdTime
	}

	updatedTime := bson.M{}
	updatedStart, _ := ctx.URLParamInt64("updatedStart")
	if updatedStart > 0 {
		updatedTime["$gte"] = time.Unix(updatedStart, 0)
	}
	updatedEnd, _ := ctx.URLParamInt64("updatedEnd")
	if updatedEnd > 0 {
		updatedTime["$lte"] = time.Unix(updatedEnd, 0)
	}
	if updatedTime["$gte"] != nil || updatedTime["$lte"] != nil {
		where["updatedAt"] = updatedTime
	}

	sort := ctx.URLParam("sort")
	if sort == "" {
		sort = "-star"
	}
	list, count, err := GetTopologies(&where, sort, pageIndex, pageCount, ctx.URLParam("count") != "0")
	if err != nil {
		ret["error"] = keys.ErrorRead
		ret["errorDetail"] = err.Error()
	}

	ret["list"] = list
	ret["count"] = count
}

// UserTopologies 获取用户拓扑图
// [query] pageIndex - 当前第几页
// [query] pageCount - 每页显示个数
// [query] count - 0，表示不统计总数，返回count为0
// [query] name - 搜索name
// [query] desc - 搜索desc
// [query] text - 搜索name和desc
// [query] createdStart
// [query] createdEnd
// [query] updatedStart
// [query] updatedEnd
// [query] deleted - 查询已删除
// [query] sort - 排序
func UserTopologies(ctx iris.Context) {
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

	name := ctx.URLParam("name")
	if name != "" {
		where["name"] = bson.M{"$regex": name, "$options": "$i"}
	}
	desc := ctx.URLParam("desc")
	if desc != "" {
		where["desc"] = bson.M{"$regex": desc, "$options": "$i"}
	}
	text := ctx.URLParam("text")
	if text != "" {
		where["$or"] = []bson.M{
			bson.M{"$regex": name, "$options": "$i"},
			bson.M{"$regex": desc, "$options": "$i"},
		}
	}

	createdTime := bson.M{}
	createdStart, _ := ctx.URLParamInt64("createdStart")
	if createdStart > 0 {
		createdTime["$gte"] = time.Unix(createdStart, 0)
	}
	createdEnd, _ := ctx.URLParamInt64("createdEnd")
	if createdEnd > 0 {
		createdTime["$lte"] = time.Unix(createdEnd, 0)
	}
	if createdTime["$gte"] != nil || createdTime["$lte"] != nil {
		where["createdAt"] = createdTime
	}

	updatedTime := bson.M{}
	updatedStart, _ := ctx.URLParamInt64("updatedStart")
	if updatedStart > 0 {
		updatedTime["$gte"] = time.Unix(updatedStart, 0)
	}
	updatedEnd, _ := ctx.URLParamInt64("updatedEnd")
	if updatedEnd > 0 {
		updatedTime["$lte"] = time.Unix(updatedEnd, 0)
	}
	if updatedTime["$gte"] != nil || updatedTime["$lte"] != nil {
		where["updatedAt"] = updatedTime
	}

	sort := ctx.URLParam("sort")
	if sort == "" {
		sort = "-updatedAt"
	}
	list, count, err := GetTopologies(&where, sort, pageIndex, pageCount, ctx.URLParam("count") != "0")
	if err != nil {
		ret["error"] = keys.ErrorRead
		ret["errorDetail"] = err.Error()
	}

	ret["list"] = list
	ret["count"] = count
}

// UserTopologyAdd 新增用户拓扑图
func UserTopologyAdd(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	data := &Topology{}
	err := ctx.ReadJSON(data)
	if err != nil {
		ret["error"] = keys.ErrorParam
		ret["errorDetail"] = err.Error()
		return
	}

	data.ID = ""
	vip := middlewares.Vip(ctx)
	err = Save(data, ctx.Values().GetString("uid"), ctx.Values().GetString("username"), vip > 0)
	if err != nil {
		ret["error"] = keys.ErrorSave
		ret["errorDetail"] = err.Error()
	}

	ret["id"] = data.ID
}

// UserTopologySave 保存修改用户拓扑图
func UserTopologySave(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	data := &Topology{}
	err := ctx.ReadJSON(data)
	if err != nil || data.ID == "" {
		ret["error"] = keys.ErrorParam
		if err != nil {
			ret["errorDetail"] = err.Error()
		}
		return
	}

	vip := middlewares.Vip(ctx)
	err = Save(data, ctx.Values().GetString("uid"), ctx.Values().GetString("username"), vip > 0)
	if err != nil {
		ret["error"] = keys.ErrorSave
		ret["errorDetail"] = err.Error()
	}
	ret["id"] = data.ID
}

// UserTopologyPatch 修改图以外的数据
func UserTopologyPatch(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	data := bson.M{}
	err := ctx.ReadJSON(&data)
	if err != nil {
		ret["error"] = keys.ErrorParam
		ret["errorDetail"] = err.Error()
		return
	}

	id := utils.String(data["id"])
	if !bson.IsObjectIdHex(id) {
		ret["error"] = keys.ErrorID
		return
	}

	err = Patch(data, ctx.Values().GetString("uid"), ctx.Values().GetString("username"))
	if err != nil {
		ret["error"] = keys.ErrorSave
		ret["errorDetail"] = err.Error()
	}
}

// UserTopologyDel 删除用户拓扑图
func UserTopologyDel(ctx iris.Context) {
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

// UserTopologyRestore 恢复删除用户拓扑图
func UserTopologyRestore(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	id := ctx.Params().Get("id")
	if !bson.IsObjectIdHex(id) {
		ret["error"] = keys.ErrorID
		return
	}
	err := Restore(id, ctx.Values().GetString("uid"))
	if err != nil {
		ret["error"] = keys.ErrorPermission
	}
}

// UserFavorites 获取用户收藏拓扑图
// [query] pageIndex - 当前第几页
// [query] pageCount - 每页显示个数
func UserFavorites(ctx iris.Context) {
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

	uid := ctx.Values().GetString("uid")
	ids, count, err := Favorites(&bson.M{"userId": uid}, pageIndex, pageCount)
	if err != nil {
		ret["error"] = keys.ErrorRead
		ret["errorDetail"] = err.Error()
		return
	}

	idList := make([]bson.ObjectId, count)
	i := 0
	for ; i < count; i++ {
		idList[i] = ids[i].ID
	}

	where := bson.M{
		"_id":    bson.M{"$in": idList},
		"userId": uid,
	}
	list, count, err := GetTopologies(&where, "", pageIndex, pageCount, ctx.URLParam("count") != "0")
	if err != nil {
		ret["warning"] = "您已收藏"
		ret["errorDetail"] = err.Error()
	}

	ret["list"] = list
	ret["count"] = count

	stars, _, _ := Stars(&where, pageIndex, pageCount)
	ret["stars"] = stars
}

// UserFavoriteAdd 收藏
func UserFavoriteAdd(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	data := &Favorite{}
	err := ctx.ReadJSON(data)
	if err != nil {
		ret["error"] = keys.ErrorParam
		ret["errorDetail"] = err.Error()
		return
	}

	if data.ID == "" {
		ret["error"] = keys.ErrorID
		return
	}

	err = FavoriteAdd(data, ctx.Values().GetString("uid"))
	if err != nil {
		ret["error"] = keys.ErrorSave
		ret["errorDetail"] = err.Error()
	}
}

// UserFavoriteDel 取消收藏
func UserFavoriteDel(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	id := ctx.Params().Get("id")
	if !bson.IsObjectIdHex(id) {
		ret["error"] = keys.ErrorID
		return
	}

	err := FavoriteDel(id, ctx.Values().GetString("uid"))
	if err != nil {
		ret["error"] = keys.ErrorSave
		ret["errorDetail"] = err.Error()
	}
}

// UserStars 获取用户点赞列表
// [query] pageIndex - 当前第几页
// [query] pageCount - 每页显示个数
func UserStars(ctx iris.Context) {
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

	uid := ctx.Values().GetString("uid")
	ids, count, err := Stars(&bson.M{"userId": uid}, pageIndex, pageCount)
	if err != nil {
		ret["error"] = keys.ErrorRead
		ret["errorDetail"] = err.Error()
		return
	}

	idList := make([]bson.ObjectId, count)
	i := 0
	for ; i < count; i++ {
		idList[i] = ids[i].ID
	}

	where := bson.M{
		"_id":    bson.M{"$in": idList},
		"userId": uid,
	}
	list, count, err := GetTopologies(&where, "", pageIndex, pageCount, ctx.URLParam("count") != "0")
	if err != nil {
		ret["error"] = keys.ErrorRead
		ret["errorDetail"] = err.Error()
	}

	ret["list"] = list
	ret["count"] = count
}

// UserStarAdd 点赞
func UserStarAdd(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	data := &Star{}
	err := ctx.ReadJSON(data)
	if err != nil {
		ret["error"] = keys.ErrorParam
		ret["errorDetail"] = err.Error()
		return
	}

	if data.ID == "" {
		ret["error"] = keys.ErrorID
		return
	}

	err = StarAdd(data, ctx.Values().GetString("uid"))
	if err != nil {
		ret["warning"] = "您已点赞"
		ret["errorDetail"] = err.Error()
	}
}

// UserStarDel 取消点赞
func UserStarDel(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	id := ctx.Params().Get("id")
	if !bson.IsObjectIdHex(id) {
		ret["error"] = keys.ErrorID
		return
	}

	err := StarDel(id, ctx.Values().GetString("uid"))
	if err != nil {
		ret["error"] = keys.ErrorSave
		ret["errorDetail"] = err.Error()
	}
}

// TopologyHistories 获取用户收藏拓扑图
// [query] fileId - 根据文件ID查找该文件历史记录
// [query] pageIndex - 当前第几页
// [query] pageCount - 每页显示个数
func TopologyHistories(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	fileID := ctx.URLParam("fileId")
	if fileID == "" || !bson.IsObjectIdHex(fileID) {
		ret["error"] = keys.ErrorParam
		return
	}

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

	list, count, err := Histories(fileID, ctx.Values().GetString("uid"), pageIndex, pageCount)
	if err != nil {
		ret["error"] = keys.ErrorRead
		ret["errorDetail"] = err.Error()
	}

	ret["list"] = list
	ret["count"] = count
}

// TopologyHistoryPatch 修改图以外的数据
func TopologyHistoryPatch(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	data := bson.M{}
	err := ctx.ReadJSON(&data)
	if err != nil {
		ret["error"] = keys.ErrorParam
		ret["errorDetail"] = err.Error()
		return
	}

	id := utils.String(data["id"])
	if !bson.IsObjectIdHex(id) {
		ret["error"] = keys.ErrorID
		return
	}

	err = HistoryPatch(data, ctx.Values().GetString("uid"), ctx.Values().GetString("username"))
	if err != nil {
		ret["error"] = keys.ErrorSave
		ret["errorDetail"] = err.Error()
	}
}

// TopologyHistoryDel 删除历史
func TopologyHistoryDel(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	id := ctx.Params().Get("id")
	if !bson.IsObjectIdHex(id) {
		ret["error"] = keys.ErrorID
		return
	}

	err := HistoryDel(id, ctx.Values().GetString("uid"))
	if err != nil {
		ret["error"] = keys.ErrorSave
		ret["errorDetail"] = err.Error()
	}
}
