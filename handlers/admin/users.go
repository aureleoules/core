package adminhandlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/backpulse/core/database"
	"github.com/backpulse/core/models"
	"github.com/backpulse/core/utils"
	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

// CreateUser : Handler for `POST /users`
func CreateUser(w http.ResponseWriter, r *http.Request) {
	/* Parse JSON into models.User */

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "error", nil)
		return
	}

	ok, errStr := database.VerifyEmail(user.Email, r)
	if !ok {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, errStr, nil)
		return
	}

	/* Check password length */
	if len(user.Password) < 8 {
		utils.RespondWithJSON(w, http.StatusUnprocessableEntity, "password_too_short", nil)
		return
	}
	if len(user.Password) > 128 {
		utils.RespondWithJSON(w, http.StatusUnprocessableEntity, "password_too_long", nil)
		return
	}

	/* Hash password */
	password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	user.Password = string(password)

	/* Insert user in Database */
	user, err = database.AddUser(user)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	tokenString, err := utils.NewJWT(user, 48)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	/* Send token in payload */
	utils.RespondWithJSON(w, http.StatusOK, "success", tokenString)

	verification, _ := database.CreateVerification(user)
	err = utils.SendVerificationMail(user.Email, verification)
	if err != nil {
		log.Println(err)
	}
	return
}

// AuthenticateUser authenticates user
func AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	/* Load config (secret) */

	var data models.User
	/* Parse json to models.User */
	_ = json.NewDecoder(r.Body).Decode(&data)

	user, err := database.GetUser(data.Email)
	if err != nil {
		/* User wasn't found */
		utils.RespondWithJSON(w, http.StatusUnauthorized, "auth_fail", nil)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		/* Password is not correct */
		utils.RespondWithJSON(w, http.StatusUnauthorized, "auth_fail", nil)
		return
	}

	/* Congratulation! You made it! Let's give a you a JWT */
	tokenString, err := utils.NewJWT(user, 48)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	/* Send token in payload */
	utils.RespondWithJSON(w, http.StatusOK, "success", tokenString)
	return
}

// GetSelfUser return self data
func GetSelfUser(w http.ResponseWriter, r *http.Request) {
	id := utils.GetUserObjectID(r)
	user, err := database.GetUserByID(id)

	log.Println(err)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	user.Password = ""
	utils.RespondWithJSON(w, http.StatusOK, "success", user)
	return
}

func RemoveSubscription(w http.ResponseWriter, r *http.Request) {
	id := utils.GetUserObjectID(r)
	user, err := database.GetUserByID(id)

	_, err = utils.StripeClient.Subscriptions.Cancel(user.ActiveSubscriptionID, nil)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	err = database.RemoveUserSubscription(user)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}

func ChargeAccount(w http.ResponseWriter, r *http.Request) {

	type chargeRequest struct {
		Name    string       `json:"name"`
		Address string       `json:"address"`
		State   string       `json:"state"`
		ZIP     string       `json:"zip"`
		City    string       `json:"city"`
		Country string       `json:"country"`
		Token   stripe.Token `json:"token"`
	}
	var data chargeRequest

	id := utils.GetUserObjectID(r)

	user, err := database.GetUserByID(id)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	/* Parse json */
	_ = json.NewDecoder(r.Body).Decode(&data)

	if data.Address == "" {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "missing_address", nil)
		return
	}
	if data.City == "" {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "missing_city", nil)
		return
	}
	if data.State == "" {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "missing_state", nil)
		return
	}
	if data.ZIP == "" {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "missing_zip", nil)
		return
	}
	if data.Country == "" {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "missing_country", nil)
		return
	}

	if data.Name == "" {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "missing_name", nil)
		return
	}

	if data.Token.ID == "" {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "missing_token", nil)
		return
	}

	planID := "plan_EPMR80fLcjkok9"

	var customer *stripe.Customer

	if user.StripeID == "" {
		/* Didn't find any custom, let's create it */
		customerParams := &stripe.CustomerParams{
			Email: stripe.String(user.Email),
			Shipping: &stripe.CustomerShippingDetailsParams{
				Address: &stripe.AddressParams{
					City:       stripe.String(data.City),
					Country:    stripe.String(data.Country),
					Line1:      stripe.String(data.Address),
					PostalCode: stripe.String(strings.Replace(data.ZIP, "_", "", -1)),
					State:      stripe.String(data.State),
				},
				Name: stripe.String(data.Name),
			},
		}
		customerParams.SetSource(data.Token.ID)

		var err error
		customer, err = utils.StripeClient.Customers.New(customerParams)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
			return
		}
	} else {
		customer, _ = utils.StripeClient.Customers.Get(user.StripeID, nil)
		if customer.Subscriptions.TotalCount > 0 {
			utils.RespondWithJSON(w, http.StatusNotAcceptable, "already_pro", nil)
			return
		}
	}

	subscription, err := utils.StripeClient.Subscriptions.New(&stripe.SubscriptionParams{
		Customer: stripe.String(customer.ID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Plan: stripe.String(planID),
			},
		},
	})

	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	err = database.SetUserPro(user.ID, customer, subscription)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}

// UpdateUser : update self data
func UpdateUser(w http.ResponseWriter, r *http.Request) {

	id := utils.GetUserObjectID(r)

	user, err := database.GetUserByID(id)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
	}

	var userData models.User
	err = json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, "error", nil)
		return
	}

	if len(userData.FullName) > 35 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "name_too_long", nil)
		return
	}

	if len(user.Country) > 100 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "country_too_long", nil)
		return
	}

	if len(user.City) > 100 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "city_too_long", nil)
		return
	}

	if len(user.Address) > 100 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "address_too_long", nil)
		return
	}
	if len(user.ZIP) > 100 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "zip_too_long", nil)
		return
	}
	if len(user.State) > 100 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "state_too_long", nil)
		return
	}

	if user.Email != userData.Email {
		//Email is now different
		//We need to check multiple things
		ok, errStr := database.VerifyEmail(userData.Email, r)
		if !ok {
			utils.RespondWithJSON(w, http.StatusNotAcceptable, errStr, nil)
			return
		}

		verification, err := database.CreateVerification(models.User{
			ID:    user.ID,
			Email: userData.Email,
		})
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
			return
		}
		err = utils.SendVerificationMail(userData.Email, verification)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
			return
		}

	}

	newUser, err := database.UpdateUser(user, userData)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	tokenString, err := utils.NewJWT(newUser, 48)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	/* Send new token in payload */
	utils.RespondWithJSON(w, http.StatusOK, "success", tokenString)
	return
}

// DeleteUser : remove user from db and also remove sheets & workbooks
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := utils.GetUserObjectID(r)

	sites, err := database.GetOwnedSites(id)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}
	if len(sites) > 0 {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "delete_sites_first", nil)
		return
	}

	err = database.RemoveUser(id)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}

// UpdateUserPassword : update user password
func UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	id := utils.GetUserObjectID(r)

	user, err := database.GetUserByID(id)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	type passwordChangeRequest struct {
		LastPassword string `json:"last_password"`
		NewPassword  string `json:"new_password"`
	}
	var request passwordChangeRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	/* Check if database's password hash is same as user request */
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.LastPassword))
	if err != nil {
		utils.RespondWithJSON(w, http.StatusNotAcceptable, "wrong_last_password", nil)
		return
	}

	/* It is */
	/* Let's check if it a valid password */
	/* Now let's hash the "new password" */

	/* Check password length */
	if len(request.NewPassword) < 8 {
		utils.RespondWithJSON(w, http.StatusUnprocessableEntity, "password_too_short", nil)
		return
	}
	if len(request.NewPassword) > 128 {
		utils.RespondWithJSON(w, http.StatusUnprocessableEntity, "password_too_long", nil)
		return
	}
	newPasswordHash, _ := bcrypt.GenerateFromPassword([]byte(request.NewPassword), 10)
	newPasswordHashString := string(newPasswordHash)

	err = database.UpdateUserPassword(id, newPasswordHashString)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "success", nil)
	return
}

// VerifyUser : verify user email
func VerifyUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strID := vars["id"]
	if bson.IsObjectIdHex(strID) {
		id := bson.ObjectIdHex(strID)
		verification, err := database.GetVerificationByID(id)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
			return
		}
		if time.Until(verification.ExpireAt) > 0 {
			//Verification didn't expire
			user, err := database.GetUserByID(verification.UserID)
			if err != nil {
				utils.RespondWithJSON(w, http.StatusNotFound, "not_found", nil)
				return
			}
			err = database.VerifyUser(user, verification)
			if err != nil {
				utils.RespondWithJSON(w, http.StatusInternalServerError, "error", nil)
				return
			}

			utils.RespondWithJSON(w, http.StatusOK, "success", nil)
			return
		}
	}
}
