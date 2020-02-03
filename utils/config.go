package utils

import (
	"encoding/json"
	"os"

	"github.com/backpulse/core/models"
)

func loadConfiguration(file string) (models.Config, error) {
	var config models.Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		return models.Config{}, err
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config, nil
}

//GetConfig returns config object
func GetConfig() models.Config {
	// Use by default ./config.json
	config, err := loadConfiguration("./config.json")
	if err != nil {
		// If it doesn't exist, take environnement variables
		config = models.Config{
			URI:           os.Getenv("MONGODB_URI"),
			Database:      os.Getenv("DATABASE"),
			Secret:        os.Getenv("SECRET"),
			GmailAddress:  os.Getenv("GMAIL_ADDRESS"),
			GmailPassword: os.Getenv("GMAIL_PASSWORD"),
			StripeKey:     os.Getenv("STRIPE_KEY"),
			BucketName:    os.Getenv("BUCKET_NAME"),
			BucketPubURL:  os.Getenv("BUCKET_PUB_URL"),
		}
		return config
	}
	return config
}
