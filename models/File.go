package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// File struct
type File struct {
	ID      bson.ObjectId `json:"id" bson:"_id"`
	OwnerID bson.ObjectId `json:"owner_id" bson:"owner_id"`
	SiteID  bson.ObjectId `json:"site_id" bson:"site_id"`

	Name string  `json:"name" bson:"name"`
	URL  string  `json:"url" bson:"url"`
	Type string  `json:"type" bson:"type"`
	Size float64 `json:"size" bson:"size"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
