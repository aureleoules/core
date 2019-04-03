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
	config, _ := loadConfiguration("./config.json")
	return config
}
