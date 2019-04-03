package models

import (
	"time"

	"github.com/backpulse/core/constants"
	"gopkg.in/mgo.v2/bson"
)

// Site struct
type Site struct {
	ID      bson.ObjectId `json:"id" bson:"_id"`
	OwnerID bson.ObjectId `json:"owner_id" bson:"owner_id"`

	DisplayName string `json:"display_name" bson:"display_name"`
	Name        string `json:"name" bson:"name"`

	Modules []constants.Module `json:"modules" bson:"modules"`

	/* Emails */
	Collaborators []Collaborator `json:"collaborators" bson:"collaborators"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`

	/* Dynamic data */
	Favorite  bool    `json:"favorite" bson:"-"`
	Role      string  `json:"role" bson:"-"`
	TotalSize float64 `json:"total_size" bson:"-"`
}

type Collaborator struct {
	Email string `json:"email" bson:"email"`
	Role  string `json:"role" bson:"role"`
}
