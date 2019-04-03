package admin

import (
	adminHandlers "github.com/backpulse/core/handlers/admin"
	"github.com/gorilla/mux"
)

func handleConstants(r *mux.Router) {
	r.HandleFunc("/constants/languages", adminHandlers.GetLanguages).Methods("GET")
}
