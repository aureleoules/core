package admin

import (
	adminHandlers "github.com/backpulse/core/handlers/admin"
	"github.com/gorilla/mux"
)

func handleProjects(r *mux.Router) {
	r.Handle("/project/{id}", ProtectedRoute(adminHandlers.GetProject)).Methods("GET")
	r.Handle("/projects/{name}", ProtectedRoute(adminHandlers.GetProjects)).Methods("GET")

	r.Handle("/projects/{name}", ProtectedRoute(adminHandlers.UpdateProject)).Methods("PUT")
	r.Handle("/project/{id}", ProtectedRoute(adminHandlers.DeleteProject)).Methods("DELETE")
}
