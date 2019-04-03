package adminhandlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/backpulse/core/database"
	"github.com/backpulse/core/models"
	"github.com/backpulse/core/utils"
	"github.com/gorilla/mux"
)

// GetContact : return contact data about site
func GetContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	contact, err := database.GetContact(site.ID)
	if err != nil {
		log.Println(err)
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "error", nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "success", contact)
	return
}

// UpdateContact : update contact data
func UpdateContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	var contactContent models.ContactContent
	/* Parse json to models.User */
	err := json.NewDecoder(r.Body).Decode(&contactContent)
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

	if len(contactContent.Name) > 60 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "name_too_long", nil)
		return
	}
	if len(contactContent.Phone) > 25 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "phone_too_long", nil)
		return
	}
	if len(contactContent.Email) > 100 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "email_too_long", nil)
		return
	}
	if len(contactContent.Address) > 150 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "address_too_long", nil)
		return
	}
	if len(contactContent.FacebookURL) > 125 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "facebook_url_too_long", nil)
		return
	}
	if len(contactContent.InstagramURL) > 125 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "instagram_url_too_long", nil)
		return
	}
	if len(contactContent.TwitterURL) > 125 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "twitter_url_too_long", nil)
		return
	}

	contactContent.SiteID = site.ID
	contactContent.OwnerID = site.OwnerID

	database.UpdateContact(site.ID, contactContent)
	utils.RespondWithJSON(w, http.StatusOK, "success", nil)

	return

}
