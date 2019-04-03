package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type VideoGroup struct {
	ID      bson.ObjectId `json:"id" bson:"_id"`
	ShortID string        `json:"short_id" bson:"short_id"`

	OwnerID bson.ObjectId `json:"owner_id" bson:"owner_id"`
	SiteID  bson.ObjectId `json:"site_id" bson:"site_id"`

	Title  string  `json:"title" bson:"title"`
	Image  string  `json:"image" bson:"image"`
	Videos []Video `json:"videos" bson:"videos"`

	Index int `json:"index" bson:"index"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
