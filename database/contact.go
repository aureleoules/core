package database

import (
	"github.com/backpulse/core/models"
	"gopkg.in/mgo.v2/bson"
)

// UpdateContact : update contact content
func UpdateContact(id bson.ObjectId, content models.ContactContent) error {
	_, err := DB.C(contactCollection).Upsert(bson.M{
		"site_id": id,
	}, bson.M{
		"$set": bson.M{
			"site_id":       id,
			"owner_id":      content.OwnerID,
			"name":          content.Name,
			"phone":         content.Phone,
			"email":         content.Email,
			"address":       content.Address,
			"facebook_url":  content.FacebookURL,
			"instagram_url": content.InstagramURL,
			"twitter_url":   content.TwitterURL,
		},
	})
	return err
}

// GetContact : return contact content
func GetContact(id bson.ObjectId) (models.ContactContent, error) {
	var contactContent models.ContactContent
	err := DB.C(contactCollection).Find(bson.M{
		"site_id": id,
	}).One(&contactContent)
	return contactContent, err
}
