package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Project struct
type Project struct {
	ID      bson.ObjectId `json:"id" bson:"_id"`
	ShortID string        `json:"short_id" bson:"short_id"`

	OwnerID bson.ObjectId `json:"owner_id" bson:"owner_id"`
	SiteID  bson.ObjectId `json:"site_id" bson:"site_id"`

	Title        string        `json:"title" bson:"title"`
	Titles       []Translation `json:"titles" bson:"titles"`
	Descriptions []Translation `json:"descriptions" bson:"descriptions"`

	URL string `json:"url" bson:"url"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
