package tool

import (
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/mgo.v2/bson"

	"topology/config"
	"topology/db/mongo"
	"topology/keys"
)

// SelectFileds 要查询的字段
var SelectFileds = bson.M{"deletedAt": false}

// List 通过查询条件获取字典列表
func List(where *bson.M, pageIndex, pageCount int, cnt bool) (list []Tool, count int, err error) {
	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	query := mongoSession.DB(config.App.Mongo.Database).C(mongo.Tool).Find(where).Sort("sort")

	if cnt {
		count, err = query.Select(bson.M{"_id": true}).Count()
	}

	if pageIndex > 0 && pageCount > 0 {
		err = query.Skip((pageIndex - 1) * pageCount).Limit(pageCount).
			All(&list)
	} else {
		err = query.All(&list)
	}
	if err != nil {
		log.Error().Caller().Err(err).Str("func", "Tool.list").Msgf("Fail to read mongo: where=%v", where)
	}

	return
}

// Save 保存，新增或修改
func Save(data *Tool, uid, username string, isOperate bool) (err error) {
	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	data.UpdatedAt = time.Now().UTC()
	data.EditorID = uid
	data.EditorName = username

	if data.ID == "" {
		data.ID = bson.NewObjectId()
		data.CreatedAt = data.UpdatedAt
		data.EditorID = uid
		data.EditorName = username
	}

	where := bson.M{"_id": data.ID}

	if isOperate {
		data.EditorName = "system"
	} else {
		where["editorId"] = uid
		where["editorName"] = bson.M{"$ne": "system"}
	}

	_, err = mongoSession.DB(config.App.Mongo.Database).C(mongo.Tool).
		Upsert(where, data)

	if err != nil {
		log.Error().Caller().Err(err).Str("func", "Tool.Save").Msgf("Fail to write mongo: data=%v", data)
	}

	return
}

// Del 删除数据
func Del(id, uid, username string, isOperate bool) (err error) {
	if !bson.IsObjectIdHex(id) {
		err = errors.New(keys.ErrorID)
		return
	}

	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	data := bson.M{
		"editorId":   uid,
		"editorName": username,
		"deletedAt":  time.Now().UTC(),
	}
	where := bson.M{"_id": bson.ObjectIdHex(id)}

	if isOperate {
		data["editorName"] = "system"
	} else {
		where["editorId"] = uid
		where["editorName"] = bson.M{"$ne": "system"}
	}

	err = mongoSession.DB(config.App.Mongo.Database).C(mongo.Tool).
		Update(where, bson.M{
			"$set": data,
		})

	if err != nil {
		log.Error().Caller().Err(err).Str("func", "Tool.Del").Msgf("Fail to write mongo:  id=%s, uid=%s", id, uid)
	}

	return
}
