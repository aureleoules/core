package utils

import (
	"net/http"
	"time"

	"github.com/backpulse/core/models"
	jwt "github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
)

//NewJWT : generates new jwt
func NewJWT(user models.User, expire int) (string, error) {
	config := GetConfig()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"id":    user.ID.Hex(),
		"exp":   time.Now().Add(time.Hour * time.Duration(expire)), /* Token expires in x hours */
	})
	/* Sign token and get string */
	tokenString, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

//GetUserObjectID : return id of jwt as type object id
func GetUserObjectID(r *http.Request) bson.ObjectId {
	decoded := GetDecodedJWT(r)
	id := decoded["id"].(string)
	return bson.ObjectIdHex(id)
}

func IsAuthorized(site models.Site, user models.User) bool {
	if site.OwnerID == user.ID {
		return true
	}
	for _, collaborator := range site.Collaborators {
		if collaborator.Email == user.Email {
			return true
		}
	}
	return false
}
