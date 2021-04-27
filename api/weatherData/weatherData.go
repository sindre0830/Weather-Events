package weatherData

// WeatherData structure stores current and predicted weather data for a location.
type WeatherData struct {
	Updated string `json:"updated"`
	Now struct {
		Air_temperature         float64 `json:"air_temperature"`
		Cloud_area_fraction     float64 `json:"cloud_area_fraction"`
		Dew_point_temperature   float64 `json:"dew_point_temperature"`
		Relative_humidity       float64 `json:"relative_humidity"`
		Wind_from_direction     float64 `json:"wind_from_direction"`
		Wind_speed              float64 `json:"wind_speed"`
		Wind_speed_of_gust      float64 `json:"wind_speed_of_gust"`
	} `json:"now"`
	Today struct {
		Summary                         string  `json:"summary"`
		Confidence                      string  `json:"confidence"`
		air_temperature_max             float64 `json:"air_temperature_max"`
		air_temperature_min             float64 `json:"air_temperature_min"`
		precipitation_amount            float64 `json:"precipitation_amount"`
		precipitation_amount_max        float64 `json:"precipitation_amount_max"`
		precipitation_amount_min        float64 `json:"precipitation_amount_min"`
		Probability_of_precipitation    float64 `json:"probability_of_precipitation"`
	} `json:"today"`
}
