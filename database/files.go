package database

import (
	"github.com/backpulse/core/models"
	"gopkg.in/mgo.v2/bson"
)

func InsertFile(file models.File) error {
	err := DB.C(filesCollection).Insert(file)
	return err
}

func GetSiteFiles(id bson.ObjectId) ([]models.File, error) {
	var files []models.File
	err := DB.C(filesCollection).Find(bson.M{
		"site_id": id,
	}).All(&files)
	return files, err
}

func DeleteFile(fileID bson.ObjectId, siteID bson.ObjectId) error {
	err := DB.C(filesCollection).Remove(bson.M{
		"_id":     fileID,
		"site_id": siteID,
	})
	return err
}

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
