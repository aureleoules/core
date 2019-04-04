package database

import (
	"time"

	"github.com/backpulse/core/models"
	"gopkg.in/mgo.v2/bson"
)

// AddVideo : add video to db
func AddVideo(video models.Video) error {
	video.UpdatedAt = time.Now()
	video.CreatedAt = time.Now()
	err := DB.C(videosCollection).Insert(video)
	return err
}

// GetVideo : Return specific video by ObjectID
func GetVideo(videoID bson.ObjectId) (models.Video, error) {
	var video models.Video
	err := DB.C(videosCollection).FindId(videoID).One(&video)
	return video, err
}

// GetGroupVideos : return array of video from a videogroup
func GetGroupVideos(id bson.ObjectId) ([]models.Video, error) {
	var videos []models.Video
	err := DB.C(videosCollection).Find(bson.M{
		"video_group_id": id,
	}).All(&videos)
	return videos, err
}

// Updatevideo : update video informations (title, content, youtube_url)
func UpdateVideo(id bson.ObjectId, video models.Video) error {
	err := DB.C(videosCollection).UpdateId(id, bson.M{
		"$set": bson.M{
			"title":       video.Title,
			"content":     video.Content,
			"youtube_url": video.YouTubeURL,
		},
	})
	return err
}

// RemoveVideo : remove video from db
func RemoveVideo(id bson.ObjectId) error {
	err := DB.C(videosCollection).RemoveId(id)
	return err
}

// UpdateVideosIndexes : update order of videos
func UpdateVideosIndexes(siteID bson.ObjectId, videos []models.Video) error {
	for _, video := range videos {
		err := DB.C(videosCollection).Update(bson.M{
			"site_id": siteID,
			"_id":     video.ID,
		}, bson.M{
			"$set": bson.M{
				"index": video.Index,
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}
