package admin

import (
	adminHandlers "github.com/backpulse/core/handlers/admin"
	"github.com/gorilla/mux"
)

func handleSites(r *mux.Router) {
	/* Sites */
	r.Handle("/sites", ProtectedRoute(adminHandlers.GetSites)).Methods("GET")
	r.Handle("/sites/{name}", ProtectedRoute(adminHandlers.GetSite)).Methods("GET")
	r.Handle("/sites/{name}", ProtectedRoute(adminHandlers.UpdateSite)).Methods("PUT")
	r.Handle("/sites/{name}", ProtectedRoute(adminHandlers.DeleteSite)).Methods("DELETE")
	r.Handle("/sites", ProtectedRoute(adminHandlers.CreateSite)).Methods("POST")

	/* Favorite */
	r.Handle("/sites/favorite/{name}", ProtectedRoute(adminHandlers.Favorite)).Methods("PUT")

	/* Modules */
	r.Handle("/sites/{name}/modules/{module}", ProtectedRoute(adminHandlers.AddModule)).Methods("POST")
	r.Handle("/sites/{name}/modules/{module}", ProtectedRoute(adminHandlers.RemoveModule)).Methods("DELETE")
	r.Handle("/sites/{name}/modules", ProtectedRoute(adminHandlers.GetSiteModules)).Methods("GET")

	/* Collaborators */
	r.Handle("/sites/{name}/collaborators", ProtectedRoute(adminHandlers.GetCollaborators)).Methods("GET")
	r.Handle("/sites/{name}/collaborators/{email}", ProtectedRoute(adminHandlers.AddCollaborator)).Methods("POST")
	r.Handle("/sites/{name}/collaborators/{email}", ProtectedRoute(adminHandlers.RemoveCollaborator)).Methods("DELETE")

	/* Transfer */
	r.Handle("/sites/{name}/transfer/{email}", ProtectedRoute(adminHandlers.TransferSite)).Methods("POST")
}
