package admin

import (
	adminHandlers "github.com/backpulse/core/handlers/admin"
	"github.com/gorilla/mux"
)

func handleTracks(r *mux.Router) {
	r.Handle("/tracks/{name}/indexes", ProtectedRoute(adminHandlers.UpdateTracksIndexes)).Methods("PUT")
	r.Handle("/tracks/{name}/{albumid}", ProtectedRoute(adminHandlers.AddTrack)).Methods("POST")

	r.Handle("/tracks/{name}/{id}", ProtectedRoute(adminHandlers.GetTrack)).Methods("GET")
	r.Handle("/tracks/{name}/{id}", ProtectedRoute(adminHandlers.DeleteTrack)).Methods("DELETE")

	r.Handle("/tracks/{name}/{id}", ProtectedRoute(adminHandlers.UpdateTrack)).Methods("PUT")
}
