package routes

import (
	"github.com/backpulse/core/routes/admin"
	"github.com/backpulse/core/routes/client"
	"github.com/gorilla/mux"
)

//NewRouter creates router with route handlers
func NewRouter() *mux.Router {
	r := mux.NewRouter()

	adminRouter := r.PathPrefix("/admin").Subrouter()
	admin.HandleAdmin(adminRouter)

	clientRouter := r.PathPrefix("/{name}").Subrouter()
	client.HandleClient(clientRouter)

	return r
}
