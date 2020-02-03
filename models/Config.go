package models

//Config config struct
type Config struct {
	URI           string
	Database      string
	Secret        string
	GmailAddress  string
	GmailPassword string
	StripeKey     string
	BucketName    string
	BucketPubURL  string
}
