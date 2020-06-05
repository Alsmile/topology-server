package tool

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Tool 左侧工具栏图标数据结构
type Tool struct {
	ID bson.ObjectId `json:"id" bson:"_id"`

	Name      string `json:"name"`
	Fullname  string `json:"fullname,omitempty" bson:"fullname,omitempty"`
	Icon      string `json:"icon"`
	Image     string `json:"image,omitempty" bson:"image,omitempty"`
	SVG       string `json:"svg,omitempty" bson:"svg,omitempty"`
	DrawFn    string `json:"drawFn,omitempty" bson:"drawFn,omitempty"`
	AnchorsFn string `json:"anchorsFn,omitempty" bson:"anchorsFn,omitempty"`
	Data      bson.M `json:"data"`

	Class string `json:"class"`
	Sort  string `json:"sort"`

	Raw      bool   `json:"raw"`
	State    int    `json:"state"`
	UserID   string `json:"userId"  bson:"userId"`
	UserName string `json:"userName"  bson:"userName"`

	EditorID   string `json:"editorId"  bson:"editorId"`
	EditorName string `json:"editorName" bson:"editorName"`

	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	DeletedAt time.Time `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
}
