package database

import (
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/backpulse/core/models"
	"github.com/stripe/stripe-go"
	"gopkg.in/mgo.v2/bson"
)

func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := DB.C(usersCollection).Find(bson.M{
		"email": email,
	}).One(&user)
	return user, err
}

func RemoveUserSubscription(user models.User) error {
	err := DB.C(usersCollection).UpdateId(user.ID, bson.M{
		"$set": bson.M{
			"active_subscription_id": "",
			"professional":           false,
		},
	})
	return err
}

func SetUserPro(id bson.ObjectId, customer *stripe.Customer, subscription *stripe.Subscription) error {
	err := DB.C(usersCollection).UpdateId(id, bson.M{
		"$set": bson.M{
			"professional":           true,
			"stripe_id":              customer.ID,
			"active_subscription_id": subscription.ID,
		},
	})
	return err
}

//UpdateUser : update user data
func UpdateUser(user models.User, updateUser models.User) (models.User, error) {
	err := DB.C(usersCollection).UpdateId(user.ID, bson.M{
		"$set": bson.M{
			"fullname":   updateUser.FullName,
			"updated_at": time.Now(),
			"address":    updateUser.Address,
			"country":    updateUser.Country,
			"zip":        updateUser.ZIP,
			"city":       updateUser.City,
			"state":      updateUser.State,
		},
	})
	if err != nil {
		return models.User{}, err
	}
	var newUser models.User
	err = DB.C(usersCollection).FindId(user.ID).One(&newUser)
	if err != nil {
		return models.User{}, err
	}
	return newUser, nil
}

//InsertEmailVerification insert email verification object in db
func InsertEmailVerification(verification models.EmailVerification) error {
	err := DB.C(emailVerificationsCollection).Insert(verification)
	return err
}

//GetVerificationByID : get email verification by id
func GetVerificationByID(id bson.ObjectId) (models.EmailVerification, error) {
	var verification models.EmailVerification
	err := DB.C(emailVerificationsCollection).FindId(id).One(&verification)
	return verification, err
}

//DeleteVerification : remove verification from db
func DeleteVerification(verification models.EmailVerification) error {
	err := DB.C(emailVerificationsCollection).RemoveId(verification.ID)
	return err
}

//VerifyUser : verify user
func VerifyUser(user models.User, verification models.EmailVerification) error {
	err := DB.C(usersCollection).UpdateId(user.ID, bson.M{
		"$set": bson.M{
			"email":          verification.Email,
			"email_verified": true,
		},
	})
	if err != nil {
		return err
	}
	err = DeleteVerification(verification)
	return err
}

func RemoveUser(id bson.ObjectId) error {
	err := DB.C(usersCollection).RemoveId(id)
	return err
}

func UpdateUserPassword(id bson.ObjectId, password string) error {
	err := DB.C(usersCollection).UpdateId(id, bson.M{
		"$set": bson.M{
			"password": password,
		},
	})
	return err
}

//AddUser inserts user into Database
func AddUser(user models.User) (models.User, error) {
	id := bson.NewObjectId()
	user.ID = id
	user.CreatedAt = time.Now()
	err := DB.C(usersCollection).Insert(&user)
	return user, err
}

//GetUser get user by username or email
func GetUser(email string) (models.User, error) {
	var user models.User
	err := DB.C(usersCollection).Find(bson.M{
		"email": email,
	}).One(&user)
	return user, err
}

//GetUserByID : returns user object with id
func GetUserByID(id bson.ObjectId) (models.User, error) {
	var user models.User
	err := DB.C(usersCollection).FindId(id).One(&user)
	return user, err
}

//IsEmailRegistered checks if email already exists in db
func IsEmailRegistered(email string) bool {
	var user models.User
	_ = DB.C(usersCollection).Find(bson.M{
		"email": email,
	}).One(&user)

	return user.Email == email
}

//CreateVerification : insert a verification status in db
func CreateVerification(user models.User) (models.EmailVerification, error) {
	id := bson.NewObjectId()
	verification := models.EmailVerification{
		UserID:   user.ID,
		Email:    user.Email,
		ExpireAt: time.Now().Add(time.Hour * time.Duration(24)),
		ID:       id,
	}
	err := InsertEmailVerification(verification)
	return verification, err
}

//VerifyEmail : check email format and whether it's been assigned already or not
func VerifyEmail(email string, r *http.Request) (bool, string) {
	/* Check email */
	if !govalidator.IsEmail(email) {
		return false, "invalid_email"
	}

	/* Check if email was used already */
	if IsEmailRegistered(email) {
		return false, "email_exists"
	}
	return true, ""
}
