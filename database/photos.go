package database

import (
	"github.com/backpulse/core/models"
	"gopkg.in/mgo.v2/bson"
)

// InsertPhoto : insert photo in db
func InsertPhoto(photo models.Photo) (models.Photo, error) {
	err := DB.C(photosCollection).Insert(photo)
	return photo, err
}

// GetGalleryPhotos return photos of gallery
func GetGalleryPhotos(id bson.ObjectId) ([]models.Photo, error) {
	var photos []models.Photo
	err := DB.C(photosCollection).Find(bson.M{
		"is_gallery": true,
		"gallery_id": id,
	}).All(&photos)
	return photos, err
}

// DeletePhotos : delete multiple photos from db
func DeletePhotos(userID bson.ObjectId, ids []bson.ObjectId) error {
	_, err := DB.C(photosCollection).RemoveAll(bson.M{
		"owner_id": userID,
		"_id": bson.M{
			"$in": ids,
		},
	})
	return err
}

func GetPhotoByID(id bson.ObjectId) (models.Photo, error) {
	var photo models.Photo
	err := DB.C(photosCollection).FindId(id).One(&photo)
	return photo, err
}

// GetPhotos : return array of photos
func GetPhotos(userID bson.ObjectId, ids []bson.ObjectId) ([]models.Photo, error) {
	var photos []models.Photo
	err := DB.C(photosCollection).Find(bson.M{
		"owner_id": userID,
		"_id": bson.M{
			"$in": ids,
		},
	}).All(&photos)
	return photos, err
}

// GetSitePhotos : return photos from site
func GetSitePhotos(id bson.ObjectId) ([]models.Photo, error) {
	var photos []models.Photo
	err := DB.C(photosCollection).Find(bson.M{
		"site_id": id,
	}).All(&photos)
	return photos, err
}

// UpdatePhotosIndexes : update order of photos
func UpdatePhotosIndexes(gallery models.Gallery, photos []models.Photo) error {
	for _, photo := range photos {
		err := DB.C(photosCollection).Update(bson.M{
			"site_id": gallery.SiteID,
			"_id":     photo.ID,
		}, bson.M{
			"$set": bson.M{
				"index": photo.Index,
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}
