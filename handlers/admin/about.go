package adminhandlers

import (
	"encoding/json"
	"net/http"

	"github.com/backpulse/core/database"
	"github.com/backpulse/core/models"
	"github.com/backpulse/core/utils"
	"github.com/gorilla/mux"
)

// GetAbout : return about content of site
func GetAbout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	aboutContent, err := database.GetAbout(site.ID)

	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "success", aboutContent)
	return
}

// UpdateAbout : update about content of site
func UpdateAbout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	var aboutContent models.AboutContent
	/* Parse json to models.User */
	err := json.NewDecoder(r.Body).Decode(&aboutContent)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "error", nil)
		return
	}

	/* Check correct owner */
	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	if len(aboutContent.Name) > 60 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "name_too_long", nil)
		return
	}

	aboutContent.SiteID = site.ID
	aboutContent.OwnerID = site.OwnerID

	database.UpdateAbout(site.ID, aboutContent)
	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}
