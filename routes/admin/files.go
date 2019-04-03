package admin

import (
	adminHandlers "github.com/backpulse/core/handlers/admin"
	"github.com/gorilla/mux"
)

func handleFiles(r *mux.Router) {
	r.Handle("/files/{name}", ProtectedRoute(adminHandlers.GetFiles)).Methods("GET")
	r.Handle("/files/{name}", ProtectedRoute(adminHandlers.UploadFile)).Methods("POST")
	r.Handle("/files/{name}/{id}/{filename}", ProtectedRoute(adminHandlers.UpdateFilename)).Methods("PUT")
	r.Handle("/files/{name}/{ids}", ProtectedRoute(adminHandlers.DeleteFiles)).Methods("DELETE")
}
