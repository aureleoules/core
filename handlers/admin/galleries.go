package adminhandlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/backpulse/core/database"
	"github.com/backpulse/core/models"
	"github.com/backpulse/core/utils"
	"github.com/gorilla/mux"
	"github.com/teris-io/shortid"
	"gopkg.in/mgo.v2/bson"
)

// DeleteGallery handler
func DeleteGallery(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	gallery, err := database.GetGallery(id)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}

	site, _ := database.GetSiteByID(gallery.SiteID)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	err = database.DeleteGallery(gallery.ID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
	}
	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}

// GetGallery return gallery
func GetGallery(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	gallery, err := database.GetGallery(id)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}
	site, _ := database.GetSiteByID(gallery.SiteID)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", gallery)
	return
}

// GetGalleries : return array of gallery
func GetGalleries(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	galleries, err := database.GetGalleries(site.ID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", galleries)
	return

}

func SetGalleryPreview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteName := vars["name"]
	galleryID := vars["galleryID"]
	photoID := vars["id"]

	site, _ := database.GetSiteByName(siteName)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	gallery, err := database.GetGallery(galleryID)
	if err != nil {
		log.Print(err)
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}

	if gallery.SiteID != site.ID {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	if !bson.IsObjectIdHex(photoID) {
		log.Print(err)
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}

	photo, err := database.GetPhotoByID(bson.ObjectIdHex(photoID))
	if err != nil {
		log.Print(err)
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}
	if photo.SiteID != site.ID || photo.GalleryID != gallery.ID {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}

	err = database.SetGalleryPreview(gallery, photo)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}

func SetDefaultGallery(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteName := vars["name"]
	galleryID := vars["id"]

	site, _ := database.GetSiteByName(siteName)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	gallery, err := database.GetGallery(galleryID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}

	if gallery.SiteID != site.ID {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	err = database.SetDefaultGallery(site, gallery)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}

// UpdateGallery handler
func UpdateGallery(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	gallery, err := database.GetGallery(id)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}
	site, _ := database.GetSiteByID(gallery.SiteID)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var galleryData models.Gallery
	/* Parse json to models.Gallery */
	err = json.NewDecoder(r.Body).Decode(&galleryData)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "error", nil)
		return
	}

	err = database.UpdateGallery(gallery.ID, galleryData)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}

// CreateGallery handler
func CreateGallery(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteName := vars["name"]
	galleryName := vars["galleryName"]

	site, _ := database.GetSiteByName(siteName)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	if len(galleryName) > 150 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "name_too_long", nil)
		return
	}

	if len(galleryName) < 1 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "name_required", nil)
		return
	}

	galleries, _ := database.GetGalleries(site.ID)

	shortID, _ := shortid.Generate()
	gallery := models.Gallery{
		ID:        bson.NewObjectId(),
		OwnerID:   site.OwnerID,
		SiteID:    site.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ShortID:   shortID,
		Index:     len(galleries),
		Titles: []models.Translation{
			{
				Content:      galleryName,
				LanguageCode: "en",
				LanguageName: "English",
			},
		},
	}

	err := database.CreateGallery(gallery)
	if err != nil {
		log.Println(err)
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "success", gallery)
	return
}

func UpdateGalleriesIndexes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteName := vars["name"]

	site, _ := database.GetSiteByName(siteName)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var galleries []models.Gallery
	/* Parse json to models.Gallery */
	err := json.NewDecoder(r.Body).Decode(&galleries)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "error", nil)
		return
	}

	err = database.UpdateGalleriesIndexes(site.ID, galleries)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "error", nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}
