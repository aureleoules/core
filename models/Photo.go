package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Photo struct
type Photo struct {
	ID      bson.ObjectId `json:"id" bson:"_id"`
	OwnerID bson.ObjectId `json:"owner_id" bson:"owner_id"`
	SiteID  bson.ObjectId `json:"site_id" bson:"site_id"`

	URL    string `json:"url" bson:"url"`
	Width  int    `json:"width" bson:"width"`
	Height int    `json:"height" bson:"height"`
	Format string `json:"format" bson:"format"`

	IsGallery bool          `json:"is_gallery" bson:"is_gallery"`
	GalleryID bson.ObjectId `json:"gallery_id" bson:"gallery_id,omitempty"`

	IsProject bool          `json:"is_project" bson:"is_project"`
	ProjectID bson.ObjectId `json:"project_id" bson:"project_id,omitempty"`
	Size      float64       `json:"size" bson:"size"`

	Index int `json:"index" bson:"index"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
