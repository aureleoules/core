package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Gallery struct
type Gallery struct {
	ID      bson.ObjectId `json:"id" bson:"_id"`
	ShortID string        `json:"short_id" bson:"short_id"`

	OwnerID bson.ObjectId `json:"owner_id" bson:"owner_id"`
	SiteID  bson.ObjectId `json:"site_id" bson:"site_id"`

	DefaultGallery bool `json:"default_gallery" bson:"default_gallery"`

	PreviewPhoto   Photo         `json:"preview_photo" bson:"-"`
	PreviewPhotoID bson.ObjectId `json:"preview_photo_id" bson:"preview_photo_id,omitempty"`

	Title        string        `json:"title" bson:"title"`
	Titles       []Translation `json:"titles" bson:"titles"`
	Descriptions []Translation `json:"descriptions" bson:"descriptions"`

	Photos []Photo `json:"photos" bson:"photos,omitempty"`

	Index int `json:"index" bson:"index"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
