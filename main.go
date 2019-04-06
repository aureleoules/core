package main

import (
	"log"
	"net/http"
	"os"

	"github.com/backpulse/core/database"
	"github.com/backpulse/core/routes"
	"github.com/backpulse/core/utils"
	"github.com/rs/cors"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	config := utils.GetConfig()

	database.Connect(config.URI, config.Database)
	utils.InitGoogleCloud()
	utils.InitStripe()

	r := routes.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Access-Control-Allow-Origin", "origin", "X-Requested-With", "Authorization", "Content-Type", "Language"},
		AllowedMethods: []string{"DELETE", "POST", "GET", "PUT"},
	})

	handler := c.Handler(r)

	var port string
	if os.Getenv("PORT") == "" {
		port = ":8000"
	} else {
		port = ":" + os.Getenv("PORT")
	}

	err := http.ListenAndServe(port, handler)
	if err != nil {
		log.Fatal(err)
	}
}
