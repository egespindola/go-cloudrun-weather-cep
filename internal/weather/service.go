package weather

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

type weatherService struct{}

func NewWeatherService() WeatherService {
	return &weatherService{}
}

func (s *weatherService) GetWeatherByCoordinates(latitude, longitude string) (*WeatherData, error) {
	url := resolveWeatherConnectorURL(latitude, longitude)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var weatherData WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return nil, err
	}

	return &weatherData, nil
}

func resolveWeatherConnectorURL(latitude, longitude string) string {
	url := os.Getenv("WEATHER_CONNECTOR_URL")
	url = strings.Replace(url, "{latitude}", latitude, 1)
	url = strings.Replace(url, "{longitude}", longitude, 1)
	return url
}
