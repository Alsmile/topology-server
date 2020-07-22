package tool

import (
	"encoding/json"
	"time"

	"topology/db/redis"

	redigo "github.com/garyburd/redigo/redis"

	"github.com/rs/zerolog/log"
	"gopkg.in/mgo.v2/bson"

	"topology/config"
	"topology/db/mongo"
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

	where := bson.M{}

	if data.ID == "" {
		data.ID = bson.NewObjectId()
		data.CreatedAt = data.UpdatedAt
		data.EditorID = uid
		data.EditorName = username

		if data.Fullname != "" && data.DrawFn != "" {
			t := new(Tool)
			err := mongoSession.DB(config.App.Mongo.Database).C(mongo.Tool).Find(
				bson.M{"fullname": data.Fullname}).Select(bson.M{"_id": true}).One(&t)

			if err == nil {
				data.ID = t.ID
			}

			where["state"] = bson.M{
				"$lt": 1,
			}

			where["fullname"] = data.Fullname
			where["_id"] = data.ID
		}
	} else {
		where["_id"] = data.ID
	}

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

// Updates 保存，新增或修改
func Updates(ids []bson.ObjectId, data *bson.M, uid, username string) (err error) {
	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	_, err = mongoSession.DB(config.App.Mongo.Database).C(mongo.Tool).
		UpdateAll(bson.M{"_id": bson.M{"$in": ids}}, data)

	if err != nil {
		log.Error().Caller().Err(err).Str("func", "Tool.Updates").Msgf("Fail to write mongo: ids=%v, data=%v", ids, data)
	}

	return
}

// Del 删除数据
func Del(ids []bson.ObjectId, uid, username string) (err error) {
	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	where := bson.M{
		"$and": []bson.M{
			bson.M{"_id": bson.M{"$in": ids}},
			bson.M{
				"$or": []bson.M{
					bson.M{"state": bson.M{"$lt": 1}},
					bson.M{"state": bson.M{"$exists": false}},
				},
			},
		},
	}

	_, err = mongoSession.DB(config.App.Mongo.Database).C(mongo.Tool).RemoveAll(where)

	if err != nil {
		log.Error().Caller().Err(err).Str("func", "Tool.Del").Msgf("Fail to write mongo:  ids=%s, uid=%s", ids, uid)
	}

	return
}

// Count 统计
func Count(group string, where *bson.M) (count []bson.M, err error) {
	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	err = mongoSession.DB(config.App.Mongo.Database).C(mongo.Tool).Pipe([]bson.M{
		bson.M{"$match": where},
		bson.M{"$group": bson.M{
			"_id":   "$" + group,
			"count": bson.M{"$sum": 1},
		}},
	}).All(&count)

	countText, err := json.Marshal(count)
	if err != nil {
		log.Error().Caller().Err(err).Str("func", "Tool.Count").Msgf("Fail to json.Marshal:  err=%v", err)
	} else {
		redisConn := redis.Pool.Get()
		defer redisConn.Close()
		_, err = redisConn.Do("SET", "topologyToolsCount", countText)
	}

	if err != nil {
		log.Error().Caller().Err(err).Str("func", "Tool.Count").Msgf("Fail to Count:  err=%v", err)
	}

	return
}

// GetCount 读取统计
func GetCount() (count []bson.M, err error) {
	redisConn := redis.Pool.Get()
	defer redisConn.Close()
	countText, err := redigo.String(redisConn.Do("GET", "topologyToolsCount"))
	if err != nil {
		Count("subClassId", &bson.M{"state": 1})
		countText, err = redigo.String(redisConn.Do("GET", "topologyToolsCount"))
		if err != nil {
			log.Error().Caller().Err(err).Str("func", "Tool.GetCount").Msgf("Fail to get count from redis:  err=%v", err)
		}
	} else {
		err = json.Unmarshal([]byte(countText), &count)
	}

	if err != nil {
		log.Error().Caller().Err(err).Str("func", "Tool.GetCount").Msgf("Fail to Unmarshal:  err=%v", err)
	}

	return
}
