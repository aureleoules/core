package database

import (
	"time"

	"github.com/backpulse/core/models"
	"gopkg.in/mgo.v2/bson"
)

// CreateGallery : insert new gallery in db
func CreateGallery(gallery models.Gallery) error {
	err := DB.C(galleriesCollection).Insert(gallery)
	return err
}

// SetDefaultGallery : set a default gallery (useful for a homepage gallery for example)
func SetDefaultGallery(site models.Site, gallery models.Gallery) error {
	_, err := DB.C(galleriesCollection).UpdateAll(bson.M{
		"site_id": site.ID,
	}, bson.M{
		"$set": bson.M{
			"default_gallery": false,
		},
	})
	if err != nil {
		return err
	}
	err = DB.C(galleriesCollection).UpdateId(gallery.ID, bson.M{
		"$set": bson.M{
			"default_gallery": true,
		},
	})
	return err
}

// UpdateGalleriesIndexes : Update galleries order
func UpdateGalleriesIndexes(siteID bson.ObjectId, galleries []models.Gallery) error {
	for _, g := range galleries {
		err := DB.C(galleriesCollection).Update(bson.M{
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

// GetGallery return specific gallery
func GetGalleryByShortID(shortID string) (models.Gallery, error) {
	var gallery models.Gallery
	err := DB.C(galleriesCollection).Find(bson.M{
		"short_id": shortID,
	}).One(&gallery)

	if err != nil {
		return models.Gallery{}, err
	}

	photos, err := GetGalleryPhotos(gallery.ID)
	if err != nil {
		return models.Gallery{}, err
	}

	gallery.Photos = photos
	gallery.PreviewPhoto, _ = GetPhotoByID(gallery.PreviewPhotoID)
	gallery.Title = getDefaultGalleryTitle(gallery)

	return gallery, err
}

// GetGallery return specific gallery
func GetGallery(id bson.ObjectId) (models.Gallery, error) {
	var gallery models.Gallery
	err := DB.C(galleriesCollection).FindId(id).One(&gallery)

	if err != nil {
		return models.Gallery{}, err
	}

	photos, err := GetGalleryPhotos(gallery.ID)
	if err != nil {
		return models.Gallery{}, err
	}

	gallery.Photos = photos
	gallery.PreviewPhoto, _ = GetPhotoByID(gallery.PreviewPhotoID)
	gallery.Title = getDefaultGalleryTitle(gallery)

	return gallery, err
}

// SetGalleryPreview : set a gallery preview image
func SetGalleryPreview(gallery models.Gallery, photo models.Photo) error {
	err := DB.C(galleriesCollection).UpdateId(gallery.ID, bson.M{
		"$set": bson.M{
			"preview_photo_id": photo.ID,
		},
	})
	return err
}

// GetDefaultGallery return default gallery
func GetDefaultGallery(id bson.ObjectId) (models.Gallery, error) {
	var gallery models.Gallery
	err := DB.C(galleriesCollection).Find(bson.M{
		"site_id":         id,
		"default_gallery": true,
	}).One(&gallery)

	if err != nil {
		return models.Gallery{}, err
	}

	photos, err := GetGalleryPhotos(gallery.ID)
	if err != nil {
		return models.Gallery{}, err
	}

	gallery.PreviewPhoto, _ = GetPhotoByID(gallery.PreviewPhotoID)

	gallery.Photos = photos

	gallery.Title = getDefaultGalleryTitle(gallery)

	return gallery, err
}

// GetGalleries return site's galleries
func GetGalleries(id bson.ObjectId) ([]models.Gallery, error) {
	var galleries []models.Gallery
	err := DB.C(galleriesCollection).Find(bson.M{
		"site_id": id,
	}).All(&galleries)

	if err != nil {
		return nil, err
	}

	for i := range galleries {
		title := getDefaultGalleryTitle(galleries[i])
		galleries[i].Title = title
		galleries[i].PreviewPhoto, _ = GetPhotoByID(galleries[i].PreviewPhotoID)
	}

	return galleries, err
}

// GetDefaultGalleryTitle : return default gallery title
func getDefaultGalleryTitle(gallery models.Gallery) string {
	for i := range gallery.Titles {
		if gallery.Titles[i].LanguageCode == "en" {
			return gallery.Titles[i].Content
		}
	}
	return gallery.Titles[0].Content
}

// UpdateGallery update gallery
func UpdateGallery(id bson.ObjectId, gallery models.Gallery) error {
	err := DB.C(galleriesCollection).UpdateId(id, bson.M{
		"$set": bson.M{
			"titles":       gallery.Titles,
			"descriptions": gallery.Descriptions,
			"updated_at":   time.Now(),
		},
	})
	return err
}

// DeleteGallery remove gallery from db
func DeleteGallery(id bson.ObjectId) error {
	err := DB.C(galleriesCollection).RemoveId(id)
	if err != nil {
		return err
	}
	_, err = DB.C(photosCollection).RemoveAll(bson.M{
		"gallery_id": id,
	})
	return err
}
