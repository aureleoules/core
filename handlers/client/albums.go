package clienthandlers

import (
	"net/http"

	"github.com/backpulse/core/database"
	"github.com/backpulse/core/utils"
	"github.com/gorilla/mux"
)

// GetAlbums : return array of album
func GetAlbums(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	site, err := database.GetSiteByName(name)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}

	albums, err := database.GetAlbums(site.ID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", albums)
	return
}

// GetAlbum : return specific album
func GetAlbum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["short_id"]

	album, err := database.GetAlbumByShortID(id)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", album)
	return
}

// GetTrack : return specific track informations
func GetTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["short_id"]

	track, err := database.GetTrackByShortID(id)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", track)
	return
}
