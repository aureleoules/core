package clienthandlers

import (
	"net/http"

	"github.com/backpulse/core/database"
	"github.com/backpulse/core/utils"
	"github.com/gorilla/mux"
)

// GetArticle : return specific article
func GetArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["short_id"]
	article, err := database.GetArticleByShortID(id)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "success", article)
	return
}

// GetArticles : return array of article
func GetArticles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	site, err := database.GetSiteByName(name)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}

	articles, err := database.GetArticles(site.ID)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", articles)
	return
}
