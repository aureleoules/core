package clienthandlers

import (
	"net/http"

	"github.com/backpulse/core/database"
	"github.com/backpulse/core/utils"
	"github.com/gorilla/mux"
)

// GetProjects return array of projects
func GetProjects(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	site, err := database.GetSiteByName(name)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}

	projects, err := database.GetProjects(site.ID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, err.Error(), nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "success", projects)
	return
}
