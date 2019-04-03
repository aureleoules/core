package admin

import (
	adminHandlers "github.com/backpulse/core/handlers/admin"
	"github.com/gorilla/mux"
)

func handleVideoGroups(r *mux.Router) {

	r.Handle("/videogroups/{name}/indexes", ProtectedRoute(adminHandlers.UpdateVideoGroupsIndexes)).Methods("PUT")
	r.Handle("/videogroups/{name}", ProtectedRoute(adminHandlers.CreateVideoGroup)).Methods("POST")

	r.Handle("/videogroups/{name}", ProtectedRoute(adminHandlers.GetVideoGroups)).Methods("GET")
	r.Handle("/videogroups/{name}/{id}", ProtectedRoute(adminHandlers.GetVideoGroup)).Methods("GET")
	r.Handle("/videogroups/{name}/{id}", ProtectedRoute(adminHandlers.UpdateVideoGroup)).Methods("PUT")

	r.Handle("/videogroups/{name}/{id}", ProtectedRoute(adminHandlers.DeleteVideoGroup)).Methods("DELETE")
}
