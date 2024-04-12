package main

import (
	"log"
	"net/http"
	"os"

	"github.com/emso-c/konzek-go-assignment/config"
	"github.com/emso-c/konzek-go-assignment/src/api"
	"github.com/emso-c/konzek-go-assignment/src/database"
	"github.com/emso-c/konzek-go-assignment/src/modules/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load secrets from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Load configuration values
	cErr := config.LoadEnv(config.DEFAULT_CONFIG_PATH)
	if cErr != nil {
		log.Fatal(cErr)
	}

	// Initialize logger
	logger, lErr := logger.InitLogger()
	if lErr != nil {
		log.Fatal(lErr)
	}
	// Only the main function should close the logger
	defer logger.Close()

	// Initialize database
	db := database.GetDatabase()
	defer db.Close()

	api.Init()
	router := api.GetRouter()

	// Start the API server
	host := os.Getenv("SERVER_HOST")
	port := os.Getenv("SERVER_PORT")
	log.Print(
		"Starting server on ",
		"http://"+host+":"+port,
	)
	hErr := http.ListenAndServe(":"+port, router)
	if hErr != nil {
		log.Fatal(hErr)
	}
}
