package adminhandlers

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/asaskevich/govalidator"
	"github.com/backpulse/core/constants"
	"github.com/backpulse/core/database"
	"github.com/backpulse/core/models"
	"github.com/backpulse/core/utils"
	"github.com/gorilla/mux"
)

// GetSites : return list of sites for a given user
func GetSites(w http.ResponseWriter, r *http.Request) {
	id := utils.GetUserObjectID(r)
	user, _ := database.GetUserByID(id)

	sites, _ := database.GetSitesOfUser(user)
	utils.RespondWithJSON(w, http.StatusOK, "success", sites)
	return
}

// GetSite : return specific site
func GetSite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	for _, c := range site.Collaborators {
		if user.Email == c.Email {
			site.Role = "collaborator"
		}
	}
	if site.Role == "" {
		site.Role = "owner"
	}

	site.TotalSize = database.GetSiteTotalSize(site)
	utils.RespondWithJSON(w, http.StatusOK, "success", site)
	return
}

func RemoveModule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	moduleName := vars["module"]
	log.Println(moduleName)

	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))
	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	err := database.RemoveModule(site, moduleName)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	exists := false
	for _, m := range site.Modules {
		if m == constants.Module(moduleName) {
			exists = true
		}
	}
	if !exists {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}

func AddModule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	moduleName := vars["module"]
	log.Println(moduleName)

	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))
	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	/* Check if module exists */
	exists := utils.CheckModuleExists(moduleName)
	if !exists {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}

	/* Check if module has not already been added */
	for _, m := range site.Modules {
		if m == constants.Module(moduleName) {
			utils.RespondWithJSON(w, http.StatusConflict, "exists", nil)
			return
		}
	}

	err := database.AddModule(site, moduleName)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return

}

func GetSiteModules(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))
	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", site.Modules)
	return
}

// Favorite : favorite a site
func Favorite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))
	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	database.SetSiteFavorite(user, site.ID)

	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}

// DeleteSite : remove all galleries, projects, photos, about, contact
func DeleteSite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	id := utils.GetUserObjectID(r)
	site, _ := database.GetSiteByName(name)
	if site.OwnerID != id {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	// Delete photos

	galleries, _ := database.GetGalleries(site.ID)
	for _, g := range galleries {
		photos, _ := database.GetGalleryPhotos(g.ID)
		utils.RemoveGoogleCloudPhotos(photos)
	}

	err := database.DeleteSite(site)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}

func GetCollaborators(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))
	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	owner, _ := database.GetUserByID(site.OwnerID)
	site.Collaborators = append(site.Collaborators, models.Collaborator{
		Email: owner.Email,
		Role:  "owner",
	})

	utils.RespondWithJSON(w, http.StatusOK, "success", site.Collaborators)
	return
}

func AddCollaborator(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	email := vars["email"]

	id := utils.GetUserObjectID(r)
	site, _ := database.GetSiteByName(name)
	if site.OwnerID != id {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	if !govalidator.IsEmail(email) {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "incorrect_email", nil)
		return
	}

	user, _ := database.GetUserByID(site.OwnerID)
	if email == user.Email {
		utils.RespondWithJSON(w, http.StatusConflict, "exists", nil)
		return
	}

	/* Check if email already exists */
	for _, c := range site.Collaborators {
		if c.Email == email {
			utils.RespondWithJSON(w, http.StatusConflict, "exists", nil)
			return
		}
	}

	owner, err := database.GetUserByID(site.OwnerID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	if len(site.Collaborators) > 0 && !owner.Professional {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "upgrade_account", nil)
		return
	}

	err = database.AddCollaborator(site.ID, models.Collaborator{
		Email: email,
		Role:  "collaborator",
	})
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return

}

func TransferSite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	email := vars["email"]

	id := utils.GetUserObjectID(r)
	site, _ := database.GetSiteByName(name)
	if site.OwnerID != id {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	exists := false
	for _, c := range site.Collaborators {
		if c.Email == email {
			exists = true
			break
		}
	}
	if !exists {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}

	isPremium := utils.HasPremiumFeatures(site, database.GetSiteTotalSize(site))

	nextOwner, err := database.GetUserByEmail(email)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}
	if isPremium && !nextOwner.Professional {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "upgrade_account", nil)
		return
	}

	lastOwner, err := database.GetUserByID(site.OwnerID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	err = database.TransferSite(site, lastOwner, nextOwner)
	if err != nil {
		log.Println(err)
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}

func RemoveCollaborator(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	email := vars["email"]

	id := utils.GetUserObjectID(r)
	site, _ := database.GetSiteByName(name)
	if site.OwnerID != id {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	if !govalidator.IsEmail(email) {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "incorrect_email", nil)
		return
	}

	user, _ := database.GetUserByID(site.OwnerID)
	if email == user.Email {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	exists := false
	for _, c := range site.Collaborators {
		if c.Email == email {
			exists = true
		}
	}
	if !exists {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "not_found", nil)
		return
	}

	err := database.RemoveCollaborator(site.ID, email)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}

// UpdateSite : update site informations
func UpdateSite(w http.ResponseWriter, r *http.Request) {
	var siteData models.Site
	/* Parse json to models.User */
	err := json.NewDecoder(r.Body).Decode(&siteData)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "error", nil)
		return
	}
	vars := mux.Vars(r)
	name := vars["name"]

	id := utils.GetUserObjectID(r)
	site, _ := database.GetSiteByName(name)
	if site.OwnerID != id {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	if site.Name != siteData.Name {
		if len(siteData.Name) < 3 {
			utils.RespondWithJSON(w, http.StatusNotAcceptable, "name_too_short", nil)
			return
		}
		if len(siteData.Name) > 30 {
			utils.RespondWithJSON(w, http.StatusNotAcceptable, "name_too_long", nil)
			return
		}

		/* lowercase letters, dashes, digits */
		var re = regexp.MustCompile(`([a-z-\d]+)`)
		var str = siteData.Name

		name := re.FindString(str)
		if name != siteData.Name {
			utils.RespondWithJSON(w, http.StatusNotAcceptable, "incorrect_characters", nil)
			return
		}

		/* Check if site doesn't already exists */
		exists := database.SiteExists(siteData.Name)
		if exists {
			utils.RespondWithJSON(w, http.StatusConflict, "site_exists", nil)
			return
		}
	}

	if len(siteData.DisplayName) > 60 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "display_name_too_long", nil)
		return
	}

	err = database.UpdateSite(site.ID, siteData)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return

}

// CreateSite : add site to db
func CreateSite(w http.ResponseWriter, r *http.Request) {
	var site models.Site
	/* Parse json to models.User */
	err := json.NewDecoder(r.Body).Decode(&site)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "error", nil)
		return
	}

	if len(site.Name) < 3 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "name_too_short", nil)
		return
	}
	if len(site.Name) > 30 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "name_too_long", nil)
		return
	}

	/* lowercase letters, dashes, digits */
	var re = regexp.MustCompile(`([a-z-\d]+)`)
	var str = site.Name

	name := re.FindString(str)
	if name != site.Name {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "incorrect_characters", nil)
		return
	}

	/* Check if site doesn't already exists */
	exists := database.SiteExists(site.Name)
	if exists {
		utils.RespondWithJSON(w, http.StatusConflict, "site_exists", nil)
		return
	}

	/* Looks good */

	id := utils.GetUserObjectID(r)

	site.OwnerID = id
	site.DisplayName = site.Name

	site, err = database.CreateSite(site)
	log.Print(err)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "success", site)
	return
}
