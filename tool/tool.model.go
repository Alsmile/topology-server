package tool

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Tool 左侧工具栏图标数据结构
type Tool struct {
	ID bson.ObjectId `json:"id" bson:"_id"`

	Name  string `json:"name"  bson:"name"`
	Icon  string `json:"icon"  bson:"icon"`
	Image string `json:"image"  bson:"image"`
	Data  bson.M `json:"data" bson:"data"`

	Class string `json:"class"  bson:"class"`

	Shared   bool   `json:"shared"`
	UserID   string `json:"userId"  bson:"userId"`
	UserName string `json:"userName"  bson:"userName"`

	EditorID   string `json:"editorId"  bson:"editorId"`
	EditorName string `json:"editorName" bson:"editorName"`

	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	DeletedAt time.Time `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
}
