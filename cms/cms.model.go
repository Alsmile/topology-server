package cms

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Cms 系统配置字段数据结构
type Cms struct {
	ID bson.ObjectId `json:"id" bson:"_id"`

	Type string `json:"type"  bson:"type"`
	Data bson.M `json:"data,omitempty" bson:"data,omitempty"`

	EditorID   string `json:"editorId"  bson:"editorId"`
	EditorName string `json:"editorName" bson:"editorName"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt,omitempty"`
	DeletedAt time.Time `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
}
