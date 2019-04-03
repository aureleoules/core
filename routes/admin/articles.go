package admin

import (
	adminHandlers "github.com/backpulse/core/handlers/admin"
	"github.com/gorilla/mux"
)

func handleArticles(r *mux.Router) {
	r.Handle("/articles/{name}", ProtectedRoute(adminHandlers.GetArticles)).Methods("GET")
	r.Handle("/articles/{name}/{id}", ProtectedRoute(adminHandlers.GetArticle)).Methods("GET")
	r.Handle("/articles/{name}", ProtectedRoute(adminHandlers.UpdateArticle)).Methods("PUT")

	r.Handle("/articles/{name}/{id}", ProtectedRoute(adminHandlers.DeleteArticle)).Methods("DELETE")
}
