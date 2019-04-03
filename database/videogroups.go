package database

import (
	"time"

	"github.com/backpulse/core/models"
	"gopkg.in/mgo.v2/bson"
)

func AddVideoGroup(videoGroup models.VideoGroup) error {
	videoGroup.UpdatedAt = time.Now()
	videoGroup.CreatedAt = time.Now()
	err := DB.C(videoGroupsCollection).Insert(videoGroup)
	return err
}

func UpdateVideoGroup(id bson.ObjectId, group models.VideoGroup) error {
	err := DB.C(videoGroupsCollection).UpdateId(id, bson.M{
		"$set": bson.M{
			"title": group.Title,
			"image": group.Image,
		},
	})
	return err
}

func GetVideoGroup(ID bson.ObjectId) (models.VideoGroup, error) {
	var videoGroup models.VideoGroup
	err := DB.C(videoGroupsCollection).FindId(ID).One(&videoGroup)

	videos, err := GetGroupVideos(videoGroup.ID)
	if err != nil {
		return models.VideoGroup{}, nil
	}
	videoGroup.Videos = videos

	return videoGroup, err
}

func GetVideoGroupByShortID(id string) (models.VideoGroup, error) {
	var videoGroup models.VideoGroup
	err := DB.C(videoGroupsCollection).Find(bson.M{
		"short_id": id,
	}).One(&videoGroup)

	videos, err := GetGroupVideos(videoGroup.ID)
	if err != nil {
		return models.VideoGroup{}, nil
	}
	videoGroup.Videos = videos

	return videoGroup, err
}

func GetVideoGroups(siteID bson.ObjectId) ([]models.VideoGroup, error) {
	var videoGroups []models.VideoGroup
	err := DB.C(videoGroupsCollection).Find(bson.M{
		"site_id": siteID,
	}).All(&videoGroups)

	return videoGroups, err
}

func RemoveVideoGroup(id bson.ObjectId) error {
	err := DB.C(videoGroupsCollection).RemoveId(id)
	return err
}

func UpdateVideoGroupsIndexes(siteID bson.ObjectId, groups []models.VideoGroup) error {
	for _, g := range groups {
		err := DB.C(videoGroupsCollection).Update(bson.M{
			"site_id": siteID,
			"_id":     g.ID,
		}, bson.M{
			"$set": bson.M{
				"index": g.Index,
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}
