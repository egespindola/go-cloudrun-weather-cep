package weather

type WeatherService interface {
	GetWeatherByCoordinates(latitude, longitude string) (*WeatherData, error)
}
