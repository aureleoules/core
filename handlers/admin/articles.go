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

// GetArticles : return array of article of site
func GetArticles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	articles, err := database.GetArticles(site.ID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", articles)
	return
}

// GetArticle : return specific article
func GetArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	id := vars["id"]

	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	article, err := database.GetArticle(bson.ObjectIdHex(id))
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}

	if article.SiteID != site.ID {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", article)
	return

}

// UpdateArticle : create/update article
func UpdateArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var article models.Article
	/* Parse json to models.Project */
	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "error", nil)
		return
	}

	if len(article.Title) < 1 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "title_required", nil)
		return
	}

	if len(article.Title) > 200 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "title_too_long", nil)
		return
	}

	if len(article.Content) < 1 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "content_required", nil)
		return
	}

	if article.ID != "" {
		a, _ := database.GetArticle(article.ID)
		if a.SiteID != site.ID {
			utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
			return
		}
	}

	article.SiteID = site.ID
	article.OwnerID = site.OwnerID
	article, err = database.UpsertArticle(article)
	if err != nil {
		log.Println(err)
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", article)
	return

}

// DeleteArticle : remove article from db
func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	id := vars["id"]

	site, _ := database.GetSiteByName(name)
	user, _ := database.GetUserByID(utils.GetUserObjectID(r))

	if !utils.IsAuthorized(site, user) {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	article, err := database.GetArticle(bson.ObjectIdHex(id))
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}

	if article.SiteID != site.ID {
		utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	err = database.RemoveArticle(article.ID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}
