package admin

import (
	adminHandlers "github.com/backpulse/core/handlers/admin"
	"github.com/gorilla/mux"
)

func handleContact(r *mux.Router) {
	r.Handle("/contact/{name}", ProtectedRoute(adminHandlers.GetContact)).Methods("GET")
	r.Handle("/contact/{name}", ProtectedRoute(adminHandlers.UpdateContact)).Methods("PUT")
}
