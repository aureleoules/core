package models

import (
	"gopkg.in/mgo.v2/bson"
)

//ContactContent struct
type ContactContent struct {
	ID bson.ObjectId `json:"id" bson:"_id"`

	SiteID  bson.ObjectId `json:"site_id" bson:"site_id"`
	OwnerID bson.ObjectId `json:"owner_id" bson:"owner_id"`

	Name    string `json:"name" bson:"name"`
	Phone   string `json:"phone" bson:"phone"`
	Email   string `json:"email" bson:"email"`
	Address string `json:"address" bson:"address"`

	FacebookURL  string `json:"facebook_url" bson:"facebook_url"`
	InstagramURL string `json:"instagram_url" bson:"instagram_url"`
	TwitterURL   string `json:"twitter_url" bson:"twitter_url"`

	CustomFields []string `json:"custom_fields" bson:"custom_fields"`
}
