package cms

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
func List(where *bson.M, pageIndex, pageCount int, cnt bool) (list []Cms, count int, err error) {
	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	query := mongoSession.DB(config.App.Mongo.Database).C(mongo.Cms).Find(where)

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
		log.Error().Caller().Err(err).Str("func", "cms.list").Msgf("Fail to read mongo: where=%v", where)
	}

	return
}

// GetCmsByID 获取指定id的内容
func GetCmsByID(id string) (*Cms, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(keys.ErrorID)
	}

	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	data := new(Cms)
	err := mongoSession.DB(config.App.Mongo.Database).C(mongo.Cms).
		Find(bson.M{"_id": bson.ObjectIdHex(id)}).Select(SelectFileds).One(&data)

	if err != nil {
		log.Error().Caller().Err(err).Str("func", "cms.GetCmsByID").Msgf("Fail to read mongo: id=%s", id)
	}

	return data, err
}

// Save 保存，新增或修改
func Save(data *Cms, uid, username string) (err error) {
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

	_, err = mongoSession.DB(config.App.Mongo.Database).C(mongo.Cms).
		Upsert(bson.M{"_id": data.ID}, data)

	if err != nil {
		log.Error().Caller().Err(err).Str("func", "cms.Save").Msgf("Fail to write mongo: data=%v", data)
	}

	return
}

// Del 删除数据
func Del(id, uid, username string) (err error) {
	if !bson.IsObjectIdHex(id) {
		err = errors.New(keys.ErrorID)
		return
	}

	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	err = mongoSession.DB(config.App.Mongo.Database).C(mongo.Cms).
		Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{
			"$set": bson.M{
				"editorId":   uid,
				"editorName": username,
				"deletedAt":  time.Now().UTC(),
			},
		})

	if err != nil {
		log.Error().Caller().Err(err).Str("func", "cms.Del").Msgf("Fail to write mongo:  id=%s, uid=%s", id, uid)
	}

	return
}
