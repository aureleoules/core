package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//EmailVerification struct
type EmailVerification struct {
	ID       bson.ObjectId `json:"id" bson:"_id"`
	UserID   bson.ObjectId `json:"user_id" bson:"user_id"`
	Email    string        `json:"email" bson:"email"`
	ExpireAt time.Time     `json:"expire_at" bson:"expire_at"`
}
