package topology

import (
	"errors"
	"time"
	"topology/config"
	"topology/db/mongo"
	"topology/keys"
	"topology/utils"

	"github.com/rs/zerolog/log"
	"gopkg.in/mgo.v2/bson"
)

// SelectFileds 要查询的字段
var SelectFileds = bson.M{"deletedAt": false}

// GetTopologyByID 获取指定id的拓扑图
func GetTopologyByID(topoID, uid string) (*Topology, error) {
	if !bson.IsObjectIdHex(topoID) {
		return nil, errors.New(keys.ErrorID)
	}

	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	topology := new(Topology)
	err := mongoSession.DB(config.App.Mongo.Database).C(mongo.Topologies).Find(
		bson.M{"_id": bson.ObjectIdHex(topoID), "$or": []bson.M{
			bson.M{"userId": uid},
			bson.M{"shared": true},
		}}).Select(SelectFileds).One(&topology)

	if err != nil {
		log.Error().Caller().Err(err).Str("func", "topology.GetTopologyByID").Msgf("Fail to read mongo: topoID=%s, uid=%s", topoID, uid)
	}

	return topology, err
}

// GetTopologies 通过查询条件获取拓扑图
func GetTopologies(where *bson.M, sort string, pageIndex, pageCount int, cnt bool) (list []Topology, count int, err error) {
	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	query := mongoSession.DB(config.App.Mongo.Database).C(mongo.Topologies).Find(where)

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
		log.Error().Caller().Err(err).Str("func", "topology.GetTopologies").Msgf("Fail to read mongo: where=%v", where)
	}

	return
}

// Save 保存，新增或修改
func Save(data *Topology, uid, username string, isHistory bool) (err error) {
	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	data.UpdatedAt = time.Now().UTC()
	data.EditorID = uid
	data.EditorName = username

	if data.ID == "" {
		data.ID = bson.NewObjectId()
		data.CreatedAt = data.UpdatedAt
		data.UserID = uid
		data.Username = username
		data.Star = 0
		data.Hot = 0
	} else {
		src, err := GetTopologyByID(data.ID.Hex(), uid)
		if err == nil {
			data.Star = src.Star
			data.Hot = src.Hot
		}
	}

	_, err = mongoSession.DB(config.App.Mongo.Database).C(mongo.Topologies).
		Upsert(bson.M{"_id": data.ID}, data)

	if err != nil {
		log.Error().Caller().Err(err).Str("func", "topology.Save").Msgf("Fail to write mongo: data=%v", data)
	}

	if isHistory {
		history(data)
	}

	return
}

// Patch 修改部分数据，不能直接修改图的数据
func Patch(data bson.M, uid, username string) (err error) {
	if !bson.IsObjectIdHex(utils.String(data["id"])) {
		err = errors.New(keys.ErrorID)
		return
	}

	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	data["updatedAt"] = time.Now().UTC()
	data["editorId"] = uid
	data["editorName"] = username

	_id := bson.ObjectIdHex(utils.String(data["id"]))
	delete(data, "id")
	delete(data, "data")
	err = mongoSession.DB(config.App.Mongo.Database).C(mongo.Topologies).
		Update(bson.M{"_id": _id, "userId": uid}, bson.M{"$set": data})

	if err != nil {
		log.Error().Caller().Err(err).Str("func", "topology.Patch").Msgf("Fail to write mongo: data=%v", data)
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

	err = mongoSession.DB(config.App.Mongo.Database).C(mongo.Topologies).
		Update(bson.M{"_id": bson.ObjectIdHex(id), "userId": uid}, bson.M{"$set": bson.M{"deletedAt": time.Now().UTC()}})

	if err != nil {
		log.Error().Caller().Err(err).Str("func", "topology.Del").Msgf("Fail to write mongo:  id=%s, uid=%s", id, uid)
	}

	return
}

// Restore 恢复删除数据
func Restore(id, uid string) (err error) {
	if !bson.IsObjectIdHex(id) {
		err = errors.New(keys.ErrorID)
		return
	}

	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	err = mongoSession.DB(config.App.Mongo.Database).C(mongo.Topologies).
		Update(bson.M{"_id": bson.ObjectIdHex(id), "userId": uid}, bson.M{"$unset": bson.M{"deletedAt": ""}})

	if err != nil {
		log.Error().Caller().Err(err).Str("func", "topology.Restore").Msgf("Fail to write mongo:  id=%s, uid=%s", id, uid)
	}

	return
}

// history 插入历史
func history(data *Topology) {
	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	data.FileID = data.ID
	data.ID = bson.NewObjectId()
	data.CreatedAt = data.UpdatedAt
	err := mongoSession.DB(config.App.Mongo.Database).C(mongo.TopologieHistories).Insert(data)
	if err != nil {
		log.Error().Caller().Err(err).Str("func", "topology.history").Msgf("Fail to write mongo: topoID=%s", data.ID)
	}

	data.ID = data.FileID
}

// Favorites 收藏列表
func Favorites(where *bson.M, pageIndex, pageCount int) (list []Favorite, count int, err error) {
	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	query := mongoSession.DB(config.App.Mongo.Database).C(mongo.Favorites).Find(where)

	count, err = query.Select(bson.M{"_id": true}).Count()

	query = query.Sort("-createdAt")
	if pageIndex > 0 && pageCount > 0 {
		err = query.Skip((pageIndex - 1) * pageCount).Limit(pageCount).
			All(&list)
	}

	if err != nil {
		log.Error().Caller().Err(err).Str("func", "topology.Favorites").Msgf("Fail to read mongo: where=%v", where)
	}

	return
}

// FavoriteAdd 收藏
func FavoriteAdd(data *Favorite, uid string) (err error) {
	if data.ID == "" {
		err = errors.New(keys.ErrorID)
		return
	}

	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	data.CreatedAt = time.Now().UTC()
	data.UserID = uid
	err = mongoSession.DB(config.App.Mongo.Database).C(mongo.Favorites).Insert(data)
	if err != nil {
		log.Error().Caller().Err(err).Str("func", "topology.FavoriteAdd.Insert").Msgf("Fail to write mongo(Favorites): data=%v", data)
		return
	}

	err = mongoSession.DB(config.App.Mongo.Database).C(mongo.Topologies).
		UpdateId(data.ID, bson.M{"$inc": bson.M{"hot": 1}})
	if err != nil {
		log.Error().Caller().Err(err).Str("func", "topology.FavoriteAdd.Inc").Msgf("Fail to write mongo(Topologies): data=%v", data)
	}

	return
}

// FavoriteDel 取消收藏
func FavoriteDel(id, uid string) (err error) {
	if !bson.IsObjectIdHex(id) {
		err = errors.New(keys.ErrorID)
		return
	}

	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	_id := bson.ObjectIdHex(id)
	err = mongoSession.DB(config.App.Mongo.Database).C(mongo.Favorites).Remove(bson.M{"_id": _id, "userId": uid})
	if err != nil {
		log.Error().Caller().Err(err).Str("func", "topology.FavoriteDel").Msgf("Fail to write mongo(Favorites): data=%v", _id)
		return
	}

	err = mongoSession.DB(config.App.Mongo.Database).C(mongo.Topologies).UpdateId(_id, bson.M{"$inc": bson.M{"hot": -1}})
	if err != nil {
		log.Error().Caller().Err(err).Str("func", "topology.FavoriteDel").Msgf("Fail to write mongo(Topologies): data=%v", _id)
	}

	return
}

// Stars 点赞列表
func Stars(where *bson.M, pageIndex, pageCount int) (list []Star, count int, err error) {
	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	query := mongoSession.DB(config.App.Mongo.Database).C(mongo.Stars).Find(where)

	count, err = query.Select(bson.M{"_id": true}).Count()

	query = query.Sort("-createdAt")
	if pageIndex > 0 && pageCount > 0 {
		err = query.Skip((pageIndex - 1) * pageCount).Limit(pageCount).
			All(&list)
	}

	if err != nil {
		log.Error().Caller().Err(err).Str("func", "topology.Stars").Msgf("Fail to read mongo: where=%v", where)
	}

	return
}

// StarAdd 点赞
func StarAdd(data *Star, uid string) (err error) {
	if data.ID == "" {
		err = errors.New(keys.ErrorID)
		return
	}

	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	data.CreatedAt = time.Now().UTC()
	data.UserID = uid
	err = mongoSession.DB(config.App.Mongo.Database).C(mongo.Stars).Insert(data)
	if err != nil {
		log.Error().Caller().Err(err).Str("func", "topology.StarAdd.Insert").Msgf("Fail to write mongo(Stars): data=%v", data)
		return
	}

	err = mongoSession.DB(config.App.Mongo.Database).C(mongo.Topologies).
		UpdateId(data.ID, bson.M{"$inc": bson.M{"star": 1}})
	if err != nil {
		log.Error().Caller().Err(err).Str("func", "topology.StarAdd.Inc").Msgf("Fail to write mongo(Topologies): data=%v", data)
	}

	return
}

// StarDel 取消点赞
func StarDel(id, uid string) (err error) {
	if !bson.IsObjectIdHex(id) {
		err = errors.New(keys.ErrorID)
		return
	}

	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	_id := bson.ObjectIdHex(id)
	err = mongoSession.DB(config.App.Mongo.Database).C(mongo.Stars).Remove(bson.M{"_id": _id, "userId": uid})
	if err != nil {
		log.Error().Caller().Err(err).Str("func", "topology.StarDel").Msgf("Fail to write mongo(Stars): data=%v", _id)
		return
	}

	err = mongoSession.DB(config.App.Mongo.Database).C(mongo.Topologies).UpdateId(_id, bson.M{"$inc": bson.M{"star": -1}})
	if err != nil {
		log.Error().Caller().Err(err).Str("func", "topology.StarDel").Msgf("Fail to write mongo(Topologies): data=%v", _id)
	}

	return
}

// GetHistory 获取拓扑图的历史记录
func GetHistory(id, fileID, uid string) (*Topology, error) {
	if !bson.IsObjectIdHex(id) || !bson.IsObjectIdHex(fileID) {
		return nil, errors.New(keys.ErrorID)
	}

	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	topology := new(Topology)
	err := mongoSession.DB(config.App.Mongo.Database).C(mongo.TopologieHistories).
		Find(bson.M{"_id": bson.ObjectIdHex(id), "fileId": bson.ObjectIdHex(fileID), "userId": uid}).
		Select(SelectFileds).One(&topology)

	if err != nil {
		log.Error().Caller().Err(err).Str("func", "topology.GetHistory").
			Msgf("Fail to read mongo: id=%s, fileID=%s, uid=%s", id, fileID, uid)
	}

	return topology, err
}

// Histories 获取指定id的历史记录
func Histories(topoID, uid string, pageIndex, pageCount int) (list []Topology, count int, err error) {
	if !bson.IsObjectIdHex(topoID) {
		err = errors.New(keys.ErrorID)
		return
	}

	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	where := bson.M{"fileId": bson.ObjectIdHex(topoID), "userId": uid}
	query := mongoSession.DB(config.App.Mongo.Database).C(mongo.TopologieHistories).Find(where)
	count, err = query.Select(bson.M{"_id": true}).Count()

	query = query.Select(SelectFileds).Sort("-createdAt")
	if pageIndex > 0 && pageCount > 0 {
		err = query.Skip((pageIndex - 1) * pageCount).Limit(pageCount).
			All(&list)
	}

	if err != nil {
		log.Error().Caller().Err(err).Str("func", "topology.Histories").Msgf("Fail to read mongo: where=%v", where)
	}

	return
}

// HistoryPatch 修改历史非图的数据
func HistoryPatch(data bson.M, uid, username string) (err error) {
	if !bson.IsObjectIdHex(utils.String(data["id"])) {
		err = errors.New(keys.ErrorID)
		return
	}

	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	data["updatedAt"] = time.Now().UTC()
	data["editorId"] = uid
	data["editorName"] = username

	_id := bson.ObjectIdHex(utils.String(data["id"]))
	delete(data, "id")
	delete(data, "data")
	delete(data, "image")
	err = mongoSession.DB(config.App.Mongo.Database).C(mongo.TopologieHistories).
		Update(bson.M{"_id": _id, "userId": uid}, bson.M{"$set": data})

	if err != nil {
		log.Error().Caller().Err(err).Str("func", "topology.HistoryPatch").Msgf("Fail to write mongo: data=%v", data)
	}

	return
}

// HistoryDel 删除历史
func HistoryDel(id, uid string) (err error) {
	if !bson.IsObjectIdHex(id) {
		err = errors.New(keys.ErrorID)
		return
	}

	mongoSession := mongo.Session.Clone()
	defer mongoSession.Close()

	err = mongoSession.DB(config.App.Mongo.Database).C(mongo.TopologieHistories).
		Remove(bson.M{"_id": bson.ObjectIdHex(id), "userId": uid})
	if err != nil {
		log.Error().Caller().Err(err).Str("func", "topology.HistoryDel").Msgf("Fail to write mongo: id=%s", id)
	}

	return
}
