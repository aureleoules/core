package admin

import (
	adminHandlers "github.com/backpulse/core/handlers/admin"
	"github.com/gorilla/mux"
)

func handleVideos(r *mux.Router) {
	r.Handle("/videos/{name}/indexes", ProtectedRoute(adminHandlers.UpdateVideosIndexes)).Methods("PUT")
	r.Handle("/videos/{name}/{groupid}", ProtectedRoute(adminHandlers.AddVideo)).Methods("POST")

	r.Handle("/videos/{name}/{id}", ProtectedRoute(adminHandlers.GetVideo)).Methods("GET")
	r.Handle("/videos/{name}/{id}", ProtectedRoute(adminHandlers.DeleteVideo)).Methods("DELETE")

	r.Handle("/videos/{name}/{id}", ProtectedRoute(adminHandlers.UpdateVideo)).Methods("PUT")
}
