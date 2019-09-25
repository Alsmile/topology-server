package images

import (
	"errors"
	"time"
	"topology/config"
	"topology/db/mongo"
	"topology/keys"

	"github.com/rs/zerolog/log"
	"gopkg.in/mgo.v2/bson"
)

// SelectFileds 要查询的字段
var SelectFileds = bson.M{"deletedAt": false}

// List 通过查询条件获取用户上传的图片
func List(where *bson.M, sort string, pageIndex, pageCount int, cnt bool) (list []Image, count int, err error) {
	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	query := mongoSession.DB(config.App.Mongo.Database).C(mongo.Images).Find(where)

	if cnt {
		count, err = query.Select(bson.M{"_id": true}).Count()
	}

	query = query.Select(SelectFileds)
	if sort != "" {
		query = query.Sort(sort)
	}

	if pageIndex > 0 && pageCount > 0 {
		err = query.Skip((pageIndex - 1) * pageCount).Limit(pageCount).
			All(&list)
	}

	if err != nil {
		log.Error().Caller().Err(err).Str("func", "images.List").Msgf("Fail to read mongo: where=%v", where)
	}

	return
}

// Put 新增
func Put(data *Image, uid, username string) (err error) {
	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	data.ID = bson.NewObjectId()
	data.CreatedAt = time.Now().UTC()
	data.UserID = uid
	data.Username = username

	err = mongoSession.DB(config.App.Mongo.Database).C(mongo.Images).Insert(data)

	if err != nil {
		log.Error().Caller().Err(err).Str("func", "images.Put").Msgf("Fail to write mongo: data=%v", data)
	}

	return
}

// Del 删除数据
func Del(id, uid string) (err error) {
	if !bson.IsObjectIdHex(id) {
		err = errors.New(keys.ErrorID)
		return
	}

	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	err = mongoSession.DB(config.App.Mongo.Database).C(mongo.Images).
		Update(bson.M{"_id": bson.ObjectIdHex(id), "userId": uid}, bson.M{"$set": bson.M{"deletedAt": time.Now().UTC()}})

	if err != nil {
		log.Error().Caller().Err(err).Str("func", "images.Del").Msgf("Fail to write mongo:  id=%s, uid=%s", id, uid)
	}

	return
}
