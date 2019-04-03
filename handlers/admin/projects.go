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

// DeleteProject : remove project from db
func DeleteProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortID := vars["id"]

	project, _ := database.GetProject(shortID)
	site, _ := database.GetSiteById(project.SiteID)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))
	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	if project.SiteID != site.ID {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	err := database.RemoveProject(shortID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}

// GetProject : return project using shortid
func GetProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortID := vars["id"]

	project, err := database.GetProject(shortID)

	site, _ := database.GetSiteById(project.SiteID)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	if project.OwnerID != site.OwnerID {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", project)
	return
}

// GetProjects : return array of projects of site
func GetProjects(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))
	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	projects, err := database.GetProjects(site.ID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", projects)
	return
}

// UpdateProject : Update or insert new project
func UpdateProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))
	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var project models.Project
	/* Parse json to models.Project */
	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "error", nil)
		return
	}

	if len(project.Titles) < 1 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "title_required", nil)
		return
	}

	if len(project.URL) > 200 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "url_too_long", nil)
		return
	}

	project.SiteID = site.ID
	project.OwnerID = site.OwnerID
	err = database.UpsertProject(project)
	if err != nil {
		log.Println(err)
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}
