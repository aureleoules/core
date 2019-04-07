package database

import (
	"time"

	"github.com/backpulse/core/models"
	"gopkg.in/mgo.v2/bson"
)

// AddTrack : create track object in db
func AddTrack(track models.Track) error {
	track.UpdatedAt = time.Now()
	track.CreatedAt = time.Now()
	err := DB.C(tracksCollection).Insert(track)
	return err
}

// GetTrack : return specific track by ObjectID
func GetTrack(trackID bson.ObjectId) (models.Track, error) {
	var track models.Track
	err := DB.C(tracksCollection).FindId(trackID).One(&track)
	return track, err
}

// GetTrack : return specific track by shortID
func GetTrackByShortID(shortID string) (models.Track, error) {
	var track models.Track
	err := DB.C(tracksCollection).Find(bson.M{
		"short_id": shortID,
	}).One(&track)
	return track, err
}

// GetAlbumTracks : return array of track for specific album
func GetAlbumTracks(id bson.ObjectId) ([]models.Track, error) {
	var tracks []models.Track
	err := DB.C(tracksCollection).Find(bson.M{
		"album_id": id,
	}).All(&tracks)
	return tracks, err
}

// UpdateTrack : update track informations (title, url, image)
func UpdateTrack(id bson.ObjectId, track models.Track) error {
	err := DB.C(tracksCollection).UpdateId(id, bson.M{
		"$set": bson.M{
			"title":   track.Title,
			"url":     track.URL,
			"image":   track.Image,
			"content": track.Content,
		},
	})
	return err
}

// RemoveTrack : delete track from db
func RemoveTrack(id bson.ObjectId) error {
	err := DB.C(tracksCollection).RemoveId(id)
	return err
}

// UpdateTracksIndexes : update order of tracks
func UpdateTracksIndexes(siteID bson.ObjectId, tracks []models.Track) error {
	for _, track := range tracks {
		err := DB.C(tracksCollection).Update(bson.M{
			"site_id": siteID,
			"_id":     track.ID,
		}, bson.M{
			"$set": bson.M{
				"index": track.Index,
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}
