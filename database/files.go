package database

import (
	"github.com/backpulse/core/models"
	"gopkg.in/mgo.v2/bson"
)

// InsertFile : create file in db
func InsertFile(file models.File) error {
	err := DB.C(filesCollection).Insert(file)
	return err
}

// GetSiteFiles : return array of file for site
func GetSiteFiles(id bson.ObjectId) ([]models.File, error) {
	var files []models.File
	err := DB.C(filesCollection).Find(bson.M{
		"site_id": id,
	}).All(&files)
	return files, err
}

// DeleteFile : remove file from db
func DeleteFile(fileID bson.ObjectId, siteID bson.ObjectId) error {
	err := DB.C(filesCollection).Remove(bson.M{
		"_id":     fileID,
		"site_id": siteID,
	})
	return err
}

// UpdateFilename : change the name of a file
func UpdateFilename(fileID bson.ObjectId, filename string, siteID bson.ObjectId) error {
	err := DB.C(filesCollection).Update(bson.M{
		"_id":     fileID,
		"site_id": siteID,
	}, bson.M{
		"$set": bson.M{
			"name": filename,
		},
	})
	return err
}
