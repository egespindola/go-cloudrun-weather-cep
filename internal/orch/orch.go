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
var ERR_INVALID_ZIPCODE = errors.New("invalid zipcode")
var ERR_ZIPCODE_NOT_FOUND = errors.New("can not find zipcode")
var ERR_INTERNAL_SERVER_ERROR = errors.New("internal server error")

func CepHandler(c *gin.Context) {
	var response TemperatureResponse

	zipcode := c.Param("zipcode")

	if !zipcodePattern.MatchString(zipcode) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": ERR_INVALID_ZIPCODE.Error()})
		return
	}

	log.Printf("zipcode: %s", zipcode)

	cepService := cep.NewCepService()

	location, err := cepService.GetLocationByCep(zipcode)
	if err != nil {
		if errors.Is(err, cep.ErrZipcodeNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": ERR_ZIPCODE_NOT_FOUND.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": ERR_INTERNAL_SERVER_ERROR.Error()})
		return
	}

	response.Celsius = weather.CurrentWeather.Temperature
	response.Fahrenheit = response.Celsius*9/5 + 32
	response.Kelvin = response.Celsius + 273.15

	c.JSON(http.StatusOK, response)
}
