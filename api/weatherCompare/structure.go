package weatherCompare

// WeatherCompare structure stores current and predicted weather data comparisons for different locations.
//
// Functionality: Handler, get
type WeatherCompare struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Location  string  `json:"location"`
	Updated   string  `json:"updated"`
	Date	  string  `json:"date"`
	Data 	  []data  `json:"data"`
}

// data structure stores weather data for a location.
type data struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Location  string  `json:"location"`
	Updated   string  `json:"updated"`
	Instant   struct {
		AirTemperature      float64 `json:"air_temperature"`
		CloudAreaFraction   float64 `json:"cloud_area_fraction"`
		DewPointTemperature float64 `json:"dew_point_temperature"`
		RelativeHumidity    float64 `json:"relative_humidity"`
		WindSpeed           float64 `json:"wind_speed"`
		WindSpeedOfGust     float64 `json:"wind_speed_of_gust"`
		PrecipitationAmount float64 `json:"precipitation_amount"`
	} `json:"instant"`
	Predicted struct {
		AirTemperatureMax          float64 `json:"air_temperature_max"`
		AirTemperatureMin          float64 `json:"air_temperature_min"`
		PrecipitationAmount        float64 `json:"precipitation_amount"`
		PrecipitationAmountMax     float64 `json:"precipitation_amount_max"`
		PrecipitationAmountMin     float64 `json:"precipitation_amount_min"`
		ProbabilityOfPrecipitation float64 `json:"probability_of_precipitation"`
	} `json:"predicted"`
}

// locationInfo structure stores all comparison locations information.
type locationInfo struct {
	Longitude float64
	Latitude  float64
	Location  string
}
