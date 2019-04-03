package admin

import (
	adminHandlers "github.com/backpulse/core/handlers/admin"
	"github.com/gorilla/mux"
)

func handleAbout(r *mux.Router) {
	r.Handle("/about/{name}", ProtectedRoute(adminHandlers.GetAbout)).Methods("GET")
	r.Handle("/about/{name}", ProtectedRoute(adminHandlers.UpdateAbout)).Methods("PUT")
}
