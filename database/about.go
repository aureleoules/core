package database

import (
	"github.com/backpulse/core/models"
	"gopkg.in/mgo.v2/bson"
)

// GetAbout : return about content of site
func GetAbout(id bson.ObjectId) (models.AboutContent, error) {
	var about models.AboutContent
	err := DB.C(aboutCollection).Find(bson.M{
		"site_id": id,
	}).One(&about)
	return about, err
}

// UpdateAbout : update about content of site
func UpdateAbout(id bson.ObjectId, about models.AboutContent) error {
	_, err := DB.C(aboutCollection).Upsert(bson.M{
		"site_id": id,
	}, bson.M{
		"$set": bson.M{
			"site_id":      id,
			"owner_id":     about.OwnerID,
			"name":         about.Name,
			"descriptions": about.Descriptions,
		},
	})
	return err
}
