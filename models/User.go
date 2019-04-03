package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//User struct
type User struct {
	ID bson.ObjectId `json:"id" bson:"_id,omitempty"`

	FullName string `json:"fullname" bson:"fullname"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`

	Country string `json:"country" bson:"country"`
	City    string `json:"city" bson:"city"`
	Address string `json:"address" bson:"address"`
	ZIP     string `json:"zip" bson:"zip"`
	State   string `json:"state" bson:"state"`

	EmailVerified bool `json:"email_verified" bson:"email_verified"`

	StripeID             string `json:"-" bson:"stripe_id"`
	ActiveSubscriptionID string `json:"-" bson:"active_subscription_id"`
	Professional         bool   `json:"professional" bson:"professional"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`

	FavoriteSites []bson.ObjectId `json:"favorite_sites" bson:"favorite_sites"`
}
