package database

import (
	"time"

	"github.com/backpulse/core/models"
	"gopkg.in/mgo.v2/bson"
)

// AddVideoGroup : add video group to db
func AddVideoGroup(videoGroup models.VideoGroup) error {
	videoGroup.UpdatedAt = time.Now()
	videoGroup.CreatedAt = time.Now()
	err := DB.C(videoGroupsCollection).Insert(videoGroup)
	return err
}

// UpdateVideoGroup : update video group informations (title, image)
func UpdateVideoGroup(id bson.ObjectId, group models.VideoGroup) error {
	err := DB.C(videoGroupsCollection).UpdateId(id, bson.M{
		"$set": bson.M{
			"title": group.Title,
			"image": group.Image,
		},
	})
	return err
}

// GetVideoGroup : return specific video group by ObjectID
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

// GetVideoGroupByShortID : return specific video group by short_id
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

// GetVideoGroups : Return array of videogroup from site
func GetVideoGroups(siteID bson.ObjectId) ([]models.VideoGroup, error) {
	var videoGroups []models.VideoGroup
	err := DB.C(videoGroupsCollection).Find(bson.M{
		"site_id": siteID,
	}).All(&videoGroups)

	return videoGroups, err
}

// RemoveVideoGroup : remove video group from db
func RemoveVideoGroup(id bson.ObjectId) error {
	err := DB.C(videoGroupsCollection).RemoveId(id)
	return err
}

// UpdateVideoGroupsIndexes : Update order of video groups
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
