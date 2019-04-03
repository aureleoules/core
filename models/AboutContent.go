package models

import (
	"gopkg.in/mgo.v2/bson"
)

//AboutContent struct
type AboutContent struct {
	ID bson.ObjectId `json:"id" bson:"_id"`

	SiteID  bson.ObjectId `json:"site_id" bson:"site_id"`
	OwnerID bson.ObjectId `json:"owner_id" bson:"owner_id"`

	Name         string        `json:"name" bson:"name"`
	Descriptions []Translation `json:"descriptions" bson:"descriptions"`
}
