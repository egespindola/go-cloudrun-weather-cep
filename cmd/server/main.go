package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/egespindola/go-cloudrun-weather-cep/internal/orch"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using OS environment variables instead")
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
