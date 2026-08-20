package orch

import (
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/egespindola/go-cloudrun-weather-cep/internal/cep"
	"github.com/egespindola/go-cloudrun-weather-cep/internal/weather"
	"github.com/gin-gonic/gin"
)

var zipcodePattern = regexp.MustCompile(`^[0-9]{8}$`)

func CepHandler(c *gin.Context) {
	var response TemperatureResponse

	zipcode := c.Param("zipcode")

	if !zipcodePattern.MatchString(zipcode) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "invalid zipcode"})
		return
	}

	log.Printf("zipcode: %s", zipcode)

	cepService := cep.NewCepService()

	location, err := cepService.GetLocationByCep(zipcode)
	if err != nil {
		if errors.Is(err, cep.ErrZipcodeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "zipcode not found"})
			return
		}

		log.Printf("error getting location for zipcode %s: %v", zipcode, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	weatherService := weather.NewWeatherService()
	weather, err := weatherService.GetWeatherByCoordinates(location.Location.Coordinates.Latitude, location.Location.Coordinates.Longitude)
	if err != nil {
		log.Printf("error getting temperature for zipcode %s: %v", zipcode, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	response.Celsius = weather.CurrentWeather.Temperature
	response.Fahrenheit = response.Celsius*9/5 + 32
	response.Kelvin = response.Celsius + 273.15

	c.JSON(http.StatusOK, response)
}
