package clienthandlers

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/backpulse/core/database"

	"github.com/backpulse/core/utils"
)

// GetAbout : return about page
func GetAbout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	site, err := database.GetSiteByName(name)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}

	about, err := database.GetAbout(site.ID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, err.Error(), nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "success", about)
	return

}
