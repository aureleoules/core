package database

import (
	"time"

	"github.com/backpulse/core/models"
	"gopkg.in/mgo.v2/bson"
)

// AddAlbum : create album object in db
func AddAlbum(album models.Album) error {
	album.UpdatedAt = time.Now()
	album.CreatedAt = time.Now()
	err := DB.C(albumsCollection).Insert(album)
	return err
}

// UpdateAlbum : update album (title, cover, description)
func UpdateAlbum(id bson.ObjectId, album models.Album) error {
	err := DB.C(albumsCollection).UpdateId(id, bson.M{
		"$set": bson.M{
			"title":       album.Title,
			"cover":       album.Cover,
			"description": album.Description,
		},
	})
	return err
}

// GetAlbum : return specific album by ObjectID
func GetAlbum(ID bson.ObjectId) (models.Album, error) {
	var album models.Album
	err := DB.C(albumsCollection).FindId(ID).One(&album)

	tracks, err := GetAlbumTracks(album.ID)
	if err != nil {
		return models.Album{}, nil
	}
	album.Tracks = tracks

	return album, err
}

// GetAlbumByShortID : return specific album by short_id
func GetAlbumByShortID(id string) (models.Album, error) {
	var album models.Album
	err := DB.C(albumsCollection).Find(bson.M{
		"short_id": id,
	}).One(&album)

	tracks, err := GetAlbumTracks(album.ID)
	if err != nil {
		return models.Album{}, nil
	}
	album.Tracks = tracks

	return album, err
}

// GetAlbums : return array of album for specific site
func GetAlbums(siteID bson.ObjectId) ([]models.Album, error) {
	var albums []models.Album
	err := DB.C(albumsCollection).Find(bson.M{
		"site_id": siteID,
	}).All(&albums)

	return albums, err
}

// RemoveAlbum : delete album object from db
func RemoveAlbum(id bson.ObjectId) error {
	err := DB.C(albumsCollection).RemoveId(id)
	return err
}

// UpdateAlbumsIndexes : change order of albums
func UpdateAlbumsIndexes(siteID bson.ObjectId, albums []models.Album) error {
	for _, album := range albums {
		err := DB.C(albumsCollection).Update(bson.M{
			"site_id": siteID,
			"_id":     album.ID,
		}, bson.M{
			"$set": bson.M{
				"index": album.Index,
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}
