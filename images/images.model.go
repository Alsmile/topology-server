package images

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Image 用户图库数据结构
type Image struct {
	ID bson.ObjectId `json:"id" bson:"_id"`

	Image string `json:"image"`

	UserID   string `json:"userId" bson:"userId"`
	Username string `json:"username" `

	CreatedAt time.Time `json:"createdAt" bson:"createdAt,omitempty"`
	DeletedAt time.Time `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
}
