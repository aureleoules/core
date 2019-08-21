package admin

import (
	adminHandlers "github.com/backpulse/core/handlers/admin"
	"github.com/gorilla/mux"
)

func handlePhotos(r *mux.Router) {
	r.Handle("/photos/{name}/create", ProtectedRoute(adminHandlers.CreatePhoto)).Methods("POST")
	r.Handle("/photos/{name}", ProtectedRoute(adminHandlers.GetPhotos)).Methods("GET")
	r.Handle("/photos/{name}/{id}", ProtectedRoute(adminHandlers.GetPhoto)).Methods("GET")
	r.Handle("/photos/{name}/{id}", ProtectedRoute(adminHandlers.UpdatePhoto)).Methods("PUT")
	r.Handle("/photos/{name}", ProtectedRoute(adminHandlers.UploadPhoto)).Methods("POST")
	r.Handle("/photos/{name}/{id}", ProtectedRoute(adminHandlers.UpdatePhotoFile)).Methods("POST")
	r.Handle("/photos/{ids}", ProtectedRoute(adminHandlers.DeletePhotos)).Methods("DELETE")

	r.Handle("/photos/{name}/{id}/indexes", ProtectedRoute(adminHandlers.UpdatePhotosIndexes)).Methods("PUT")
}
