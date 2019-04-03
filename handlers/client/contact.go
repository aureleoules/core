package clienthandlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/backpulse/core/database"

	"github.com/backpulse/core/utils"
)

// GetContact : return contact page
func GetContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	site, err := database.GetSiteByName(name)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}

	contact, err := database.GetContact(site.ID)
	if err != nil {
		log.Println(err)
		utils.RespondWithJSON(w, http.StatusNotFound, err.Error(), nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "success", contact)
	return
}
