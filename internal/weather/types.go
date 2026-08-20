package weather

type WeatherData struct {
	Latitude             float64             `json:"latitude"`
	Longitude            float64             `json:"longitude"`
	GenerationtimeMs     float64             `json:"generationtime_ms"`
	UtcOffsetSeconds     int                 `json:"utc_offset_seconds"`
	Timezone             string              `json:"timezone"`
	Timezoneabbreviation string              `json:"timezone_abbreviation"`
	Elevation            float64             `json:"elevation"`
	CurrentWeatherUnits  CurrentWeatherUnits `json:"current_weather_units"`
	CurrentWeather       CurrentWeather      `json:"current_weather"`
}

type CurrentWeatherUnits struct {
	Time        string `json:"time"`
	Interval    string `json:"interval"`
	Temperature string `json:"temperature"`
	IsDay       string `json:"is_day"`
	WeatherCode string `json:"weather_code"`
	WindSpeed   string `json:"wind_speed"`
}

type CurrentWeather struct {
	Time          string  `json:"time"`
	Interval      int     `json:"interval"`
	Temperature   float64 `json:"temperature"`
	WindSpeed     float64 `json:"wind_speed"`
	WindDirection float64 `json:"wind_direction"`
	IsDay         int     `json:"is_day"`
	WeatherCode   int     `json:"weather_code"`
}
