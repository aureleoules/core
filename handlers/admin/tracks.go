package adminhandlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/backpulse/core/database"
	"github.com/backpulse/core/models"
	"github.com/backpulse/core/utils"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// AddTrack : add track to album
func AddTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	albumid := vars["albumid"]

	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	album, err := database.GetAlbum(bson.ObjectIdHex(albumid))
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}

	if album.SiteID != site.ID {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var track models.Track
	/* Parse json to models.Track */
	err = json.NewDecoder(r.Body).Decode(&track)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "error", nil)
		return
	}

	track.AlbumID = bson.ObjectIdHex(albumid)

	if track.AlbumID == "" {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "album_id_required", nil)
		return
	}

	track.SiteID = site.ID
	track.OwnerID = site.OwnerID

	track.ID = bson.NewObjectId()

	tracks, _ := database.GetAlbumTracks(track.AlbumID)
	track.Index = len(tracks)

	err = database.AddTrack(track)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}

// UpdateTrack : update track informations
func UpdateTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	id := vars["id"]

	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var track models.Track
	/* Parse json to models.Track */
	err := json.NewDecoder(r.Body).Decode(&track)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "error", nil)
		return
	}

	err = database.UpdateTrack(bson.ObjectIdHex(id), track)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}

// DeleteTrack : remove track from album
func DeleteTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	id := vars["id"]

	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	track, err := database.GetTrack(bson.ObjectIdHex(id))
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}

	err = database.RemoveTrack(track.ID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}

// GetTrack : get specific track
func GetTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	id := vars["id"]

	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	track, err := database.GetTrack(bson.ObjectIdHex(id))
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", track)
	return
}

// UpdateTracksIndexes : update tracks order
func UpdateTracksIndexes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteName := vars["name"]

	site, _ := database.GetSiteByName(siteName)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var tracks []models.Track
	/* Parse json to models.Gallery */
	err := json.NewDecoder(r.Body).Decode(&tracks)
	if err != nil {
		log.Print(err)
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "error", nil)
		return
	}

	err = database.UpdateTracksIndexes(site.ID, tracks)
	if err != nil {
		log.Print(err)
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "error", nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}
