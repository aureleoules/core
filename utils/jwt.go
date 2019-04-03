package utils

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

//GetDecodedJWT Returns decoded JWT object
func GetDecodedJWT(r *http.Request) jwt.MapClaims {
	return r.Context().Value("user").(*jwt.Token).Claims.(jwt.MapClaims)
}
