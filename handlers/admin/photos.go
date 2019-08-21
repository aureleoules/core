package adminhandlers

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/backpulse/core/database"
	"github.com/backpulse/core/models"
	"github.com/backpulse/core/utils"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// UpdatePhotosIndexes : change order of photos
func UpdatePhotosIndexes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteName := vars["name"]
	galleryID := vars["id"]

	site, _ := database.GetSiteByName(siteName)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	gallery, err := database.GetGallery(bson.ObjectIdHex(galleryID))
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}

	if gallery.SiteID != site.ID {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var photos []models.Photo
	/* Parse json to models.Gallery */
	err = json.NewDecoder(r.Body).Decode(&photos)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "error", nil)
		return
	}

	err = database.UpdatePhotosIndexes(gallery, photos)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}

// DeletePhotos : remove photos from db & g cloud
func DeletePhotos(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ids := vars["ids"]
	stringIDsArray := strings.Split(ids, ",")
	var array []bson.ObjectId
	log.Println(stringIDsArray)
	for _, id := range stringIDsArray {
		if bson.IsObjectIdHex(id) {
			array = append(array, bson.ObjectIdHex(id))
		}
	}

	log.Println(array)

	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	photos, err := database.GetPhotos(user.ID, array)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	var site models.Site
	for _, photo := range photos {
		if photo.SiteID != site.ID {
			site, _ = database.GetSiteByID(photo.SiteID)
		}
		if !utils.IsAuthorized(site, user) {
			utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
			return
		}
	}

	err = database.DeletePhotos(user.ID, array)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}
	err = utils.RemoveGoogleCloudPhotos(photos)
	// if err != nil {
	// 	utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
	// 	return
	// }

	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}

// GetPhotos : return photos
func GetPhotos(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteName := vars["name"]

	site, _ := database.GetSiteByName(siteName)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	photos, err := database.GetSitePhotos(site.ID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "success", photos)
	return
}

func UpdatePhotoFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	photoID := bson.ObjectIdHex(vars["id"])
	name := vars["name"]
	/* Get file from Client */
	image, header, err := r.FormFile("image")
	if err != nil {
		log.Println(err)
		utils.RespondWithJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	defer image.Close()

	if !strings.HasPrefix(header.Header.Get("Content-Type"), "image") {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "not_an_image", nil)
		return
	}
	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	photo := models.Photo{
		Size:   math.Round(float64(header.Size)/10000) / 100,
		Format: header.Header.Get("Content-Type"),
	}

	id, err := utils.UploadFile(image, header.Filename)
	if err != nil {
		log.Println(err)
		utils.RespondWithJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	config := utils.GetConfig()
	url := "https://" + config.BucketName + "/" + id.Hex()

	err = database.UpdatePhotoURL(photoID, url)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", photo)
	return

}

// UploadPhoto handler
func UploadPhoto(w http.ResponseWriter, r *http.Request) {
	/* Get file from Client */
	image, header, err := r.FormFile("image")
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	defer image.Close()

	vars := mux.Vars(r)
	name := vars["name"]

	if !strings.HasPrefix(header.Header.Get("Content-Type"), "image") {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "not_an_image", nil)
		return
	}
	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	photo := models.Photo{
		Size:      math.Round(float64(header.Size)/10000) / 100,
		CreatedAt: time.Now(),
		OwnerID:   site.OwnerID,
		Format:    header.Header.Get("Content-Type"),
		SiteID:    site.ID,
	}

	if r.FormValue("is_gallery") == "true" {
		// This is a gallery photos
		galleryID := r.FormValue("gallery_id")
		photo.IsGallery = true
		g, err := database.GetGallery(bson.ObjectIdHex(galleryID))
		if err != nil {
			utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
			return
		}
		if g.OwnerID != site.OwnerID {
			utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
			return
		}
		photo.GalleryID = &g.ID
		photo.Title = r.FormValue("title")
		photo.Content = r.FormValue("content")

	} else if r.FormValue("is_project") == "true" {
		// This is a project photo
		photo.IsProject = true
		projectID := r.FormValue("project_id")

		p, err := database.GetProject(bson.ObjectIdHex(projectID))
		if err != nil {
			utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
			return
		}
		if p.OwnerID != site.OwnerID {
			utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
			return
		}
		photo.ProjectID = p.ID
	} else {
		utils.RespondWithJSON(w, http.StatusBadRequest, "not_acceptable", nil)
		return
	}

	id, err := utils.UploadFile(image, header.Filename)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	photo.ID = id
	config := utils.GetConfig()
	photo.URL = "https://" + config.BucketName + "/" + id.Hex()

	photo, err = database.InsertPhoto(photo)
	if err != nil {
		log.Println(err)
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "success", photo)
	return
}

func CreatePhoto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteName := vars["name"]

	site, _ := database.GetSiteByName(siteName)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		log.Println("unauthorized")
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var photo models.Photo
	/* Parse json to models.Project */
	err := json.NewDecoder(r.Body).Decode(&photo)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "error", nil)
		return
	}

	photo.OwnerID = site.OwnerID
	photo.SiteID = site.ID

	log.Println(photo)

	photo, err = database.CreatePhoto(photo)
	if err != nil {
		log.Println(err)
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "success", photo)
	return

}

func GetPhoto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteName := vars["name"]
	photoID := vars["id"]

	site, _ := database.GetSiteByName(siteName)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	photo, err := database.GetPhotoByID(bson.ObjectIdHex(photoID))
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	if photo.GalleryID != nil {
		gallery, _ := database.GetGallery(*photo.GalleryID)
		photo.GalleryName = gallery.Title
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", photo)
	return
}

func UpdatePhoto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteName := vars["name"]
	photoID := vars["id"]

	site, _ := database.GetSiteByName(siteName)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var photo models.Photo
	/* Parse json to models.Project */
	err := json.NewDecoder(r.Body).Decode(&photo)
	if err != nil {
		log.Println(err)
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "error", nil)
		return
	}

	err = database.UpdatePhoto(bson.ObjectIdHex(photoID), photo)
	if err != nil {
		log.Println(err)
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}
}
