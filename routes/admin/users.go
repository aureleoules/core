package admin

import (
	adminHandlers "github.com/backpulse/core/handlers/admin"
	"github.com/gorilla/mux"
)

func handleUsers(r *mux.Router) {
	r.HandleFunc("/users", adminHandlers.CreateUser).Methods("POST")
	r.HandleFunc("/users/authenticate", adminHandlers.AuthenticateUser).Methods("POST")
	r.HandleFunc("/users/verify/{id}", adminHandlers.VerifyUser).Methods("POST")

	r.Handle("/user", ProtectedRoute(adminHandlers.DeleteUser)).Methods("DELETE")
	r.Handle("/user/password", ProtectedRoute(adminHandlers.UpdateUserPassword)).Methods("PUT")
	r.Handle("/profile", ProtectedRoute(adminHandlers.UpdateUser)).Methods("PUT")

	r.Handle("/account/charge", ProtectedRoute(adminHandlers.ChargeAccount)).Methods("POST")
	r.Handle("/account/subscription", ProtectedRoute(adminHandlers.RemoveSubscription)).Methods("DELETE")
	r.Handle("/user", ProtectedRoute(adminHandlers.GetSelfUser)).Methods("GET")
}
