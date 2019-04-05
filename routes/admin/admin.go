package admin

import (
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/backpulse/core/database"
	"github.com/backpulse/core/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

var jwtMiddleware *jwtmiddleware.JWTMiddleware

// HandleAdmin : setup all admin routes
func HandleAdmin(r *mux.Router) *mux.Router {
	config := utils.GetConfig()

	jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Secret), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	handleUsers(r)
	handleSites(r)
	handleContact(r)
	handleConstants(r)
	handleProjects(r)
	handleAbout(r)
	handlePhotos(r)
	handleGalleries(r)
	handleArticles(r)
	handleVideos(r)
	handleVideoGroups(r)
	handleFiles(r)
	handleAlbums(r)
	handleTracks(r)

	return r
}

// ProtectedRoute : returns a JWT protected route handler
func ProtectedRoute(next func(w http.ResponseWriter, r *http.Request)) *negroni.Negroni {
	return negroni.New(negroni.HandlerFunc(jwtMiddleware.HandlerWithNext), negroni.WrapFunc(func(w http.ResponseWriter, r *http.Request) {
		id := utils.GetUserObjectID(r)
		_, err := database.GetUserByID(id)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusUnauthorized, "unauthorized", nil)
			return
		}
		next(w, r)
	}))
}
