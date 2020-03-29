package adminhandlers

import (
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

// GetFiles : Return array of files for site
func GetFiles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	files, err := database.GetSiteFiles(site.ID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "success", files)
	return
}

// UpdateFilename : update filename on gcloud
func UpdateFilename(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	id := vars["id"]
	filename := vars["filename"]

	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	if len(filename) < 1 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "name_required", nil)
		return
	}

	err := database.UpdateFilename(bson.ObjectIdHex(id), filename, site.ID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}

	err = utils.UpdateFilename(bson.ObjectIdHex(id), filename)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}

// UploadFile : upload file to gcloud
func UploadFile(w http.ResponseWriter, r *http.Request) {
	/* Get file from Client */
	file, header, err := r.FormFile("file")
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	defer file.Close()

	vars := mux.Vars(r)
	name := vars["name"]

	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	id, err := utils.UploadFile(file, header.Filename)
	if err != nil {
		log.Println(err)
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	config := utils.GetConfig()

	fileObject := models.File{
		ID:        id,
		Name:      header.Filename,
		OwnerID:   site.OwnerID,
		SiteID:    site.ID,
		Size:      math.Round(float64(header.Size)/10000) / 100,
		CreatedAt: time.Now(),
		Type:      header.Header.Get("Content-Type"),
		URL:       config.BucketPubURL + "/" + id.Hex(),
	}

	err = database.InsertFile(fileObject)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "success", fileObject)
	return
}

// DeleteFiles : remove files from db & gcloud
func DeleteFiles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	fileIds := vars["ids"]

	stringIDsArray := strings.Split(fileIds, ",")
	var ids []bson.ObjectId

	for _, id := range stringIDsArray {
		if bson.IsObjectIdHex(id) {
			ids = append(ids, bson.ObjectIdHex(id))
		}
	}

	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	for _, id := range ids {
		err := database.DeleteFile(id, site.ID)
		if err != nil {
			log.Println(err)
			utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
			return
		}
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return

}
