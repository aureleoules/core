package admin

import (
	adminHandlers "github.com/backpulse/core/handlers/admin"
	"github.com/gorilla/mux"
)

func handleAlbums(r *mux.Router) {

	r.Handle("/albums/{name}/indexes", ProtectedRoute(adminHandlers.UpdateAlbumsIndexes)).Methods("PUT")
	r.Handle("/albums/{name}", ProtectedRoute(adminHandlers.CreateAlbum)).Methods("POST")

	r.Handle("/albums/{name}", ProtectedRoute(adminHandlers.GetAlbums)).Methods("GET")
	r.Handle("/albums/{name}/{id}", ProtectedRoute(adminHandlers.GetAlbum)).Methods("GET")
	r.Handle("/albums/{name}/{id}", ProtectedRoute(adminHandlers.UpdateAlbum)).Methods("PUT")

	r.Handle("/albums/{name}/{id}", ProtectedRoute(adminHandlers.DeleteAlbum)).Methods("DELETE")
}
