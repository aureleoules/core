package admin

import (
	adminHandlers "github.com/backpulse/core/handlers/admin"
	"github.com/gorilla/mux"
)

func handleGalleries(r *mux.Router) {
	r.Handle("/gallery/{id}", ProtectedRoute(adminHandlers.GetGallery)).Methods("GET")
	r.Handle("/gallery/{id}", ProtectedRoute(adminHandlers.DeleteGallery)).Methods("DELETE")
	r.Handle("/gallery/{id}", ProtectedRoute(adminHandlers.UpdateGallery)).Methods("PUT")

	r.Handle("/galleries/{name}/{galleryID}/preview/{id}", ProtectedRoute(adminHandlers.SetGalleryPreview)).Methods("PUT")
	r.Handle("/galleries/{name}/default/{id}", ProtectedRoute(adminHandlers.SetDefaultGallery)).Methods("PUT")
	r.Handle("/galleries/{name}/indexes", ProtectedRoute(adminHandlers.UpdateGalleriesIndexes)).Methods("PUT")
	r.Handle("/galleries/{name}/{galleryName}", ProtectedRoute(adminHandlers.CreateGallery)).Methods("POST")
	r.Handle("/galleries/{name}", ProtectedRoute(adminHandlers.GetGalleries)).Methods("GET")
}
