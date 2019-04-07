package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/backpulse/core/database"
	"github.com/backpulse/core/routes"
	"github.com/backpulse/core/utils"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

var (
	envFile string
)

func init() {
	// env file and will load them into ENV for this process.
	// will not overload value that defined in ENV had present.
	flag.StringVar(&envFile, "env", ".env", "env config file")
}

func main() {
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	godotenv.Load(envFile)
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
