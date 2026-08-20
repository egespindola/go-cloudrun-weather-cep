package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/egespindola/go-cloudrun-weather-cep/internal/orch"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file, using OS environment variables instead")
	}

	router := gin.Default()
	router.GET("/cep/:zipcode", orch.CepHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3099"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}

	log.Printf("Server running on port %s", port)
}
