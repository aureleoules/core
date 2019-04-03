package client

import (
	"net/http"

	clientHandlers "github.com/backpulse/core/handlers/client"
	"github.com/backpulse/core/utils"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// HandleClient : handle client routes
func HandleClient(r *mux.Router) {

	greetings := func(w http.ResponseWriter, r *http.Request) {
		utils.RespondWithJSON(w, http.StatusOK, "Welcome to the API", bson.M{
			"wrapper": "https://github.com/backpulse/wrapper",
		})
	}

	r.HandleFunc("", greetings).Methods("GET")
	r.HandleFunc("/", greetings).Methods("GET")

	r.HandleFunc("/contact", clientHandlers.GetContact).Methods("GET")
	r.HandleFunc("/about", clientHandlers.GetAbout).Methods("GET")
	r.HandleFunc("/galleries/default", clientHandlers.GetDefaultGallery).Methods("GET")
	r.HandleFunc("/galleries", clientHandlers.GetGalleries).Methods("GET")
	r.HandleFunc("/gallery/{id}", clientHandlers.GetGallery).Methods("GET")
	r.HandleFunc("/projects", clientHandlers.GetProjects).Methods("GET")
	r.HandleFunc("/articles", clientHandlers.GetArticles).Methods("GET")
	r.HandleFunc("/articles/{id}", clientHandlers.GetArticle).Methods("GET")

	r.HandleFunc("/videogroups", clientHandlers.GetVideoGroups).Methods("GET")
	r.HandleFunc("/videogroups/{id}", clientHandlers.GetVideoGroup).Methods("GET")
	r.HandleFunc("/videos/{id}", clientHandlers.GetVideo).Methods("GET")

}
