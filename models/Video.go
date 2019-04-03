package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Video struct {
	ID bson.ObjectId `json:"id" bson:"_id"`

	OwnerID bson.ObjectId `json:"owner_id" bson:"owner_id"`
	SiteID  bson.ObjectId `json:"site_id" bson:"site_id"`

	Title      string `json:"title" bson:"title"`
	Content    string `json:"content" bson:"content"`
	YouTubeURL string `json:"youtube_url" bson:"youtube_url"`

	VideoGroupID bson.ObjectId `json:"video_group_id" bson:"video_group_id"`
	Index        int           `json:"index" bson:"index"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
