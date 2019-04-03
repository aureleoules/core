package utils

import (
	"encoding/json"
	"net/http"

	"github.com/backpulse/core/models"
	gomail "gopkg.in/gomail.v2"
)

//Response Struct
type Response struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

//RespondWithJSON : Respond with JSON
func RespondWithJSON(w http.ResponseWriter, code int, message string, payload interface{}) {
	var status string
	if code >= 200 && code <= 299 {
		status = "success"
	} else {
		status = "error"
	}

	response, _ := json.MarshalIndent(Response{
		Status:  status,
		Code:    code,
		Message: message,
		Payload: payload,
	}, "", "    ")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

//SendVerificationMail : send verification email to user
func SendVerificationMail(email string, verification models.EmailVerification) error {
	config := GetConfig()
	m := gomail.NewMessage()

	m.SetHeader("From", "no-reply@backpulse.io")
	m.SetHeader("To", email)

	m.SetHeader("Subject", "Please verify your email address")

	link := "https://www.backpulse.io/verify/" + verification.ID.Hex()
	linkAsATag := "<a href=\"" + link + "\">" + link + "</a>"

	m.SetBody("text/html", `Please click the following link to confirm that <b>`+email+`</b> is your email address.<br/>`+linkAsATag+`<br/><br/><i>Thanks for using <b>Backpulse</b>!`)

	d := gomail.NewDialer("smtp.gmail.com", 465, config.GmailAddress, config.GmailPassword)

	err := d.DialAndSend(m)
	return err
	return nil
}
