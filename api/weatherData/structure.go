package weatherData

// WeatherData structure stores current and predicted weather data for the next 9 days.
//
// Functionality: Handler, get, readData
type WeatherData struct {
	Updated    string                `json:"updated"`
	Timeseries map[string]Timeseries `json:"timeseries"`
}

// Timeseries stores current and predicted weather data for a day
type Timeseries struct {
	Instant struct {
		AirTemperature      float64 `json:"air_temperature"`
		CloudAreaFraction   float64 `json:"cloud_area_fraction"`
		DewPointTemperature float64 `json:"dew_point_temperature"`
		RelativeHumidity    float64 `json:"relative_humidity"`
		WindFromDirection   float64 `json:"wind_from_direction"`
		WindSpeed           float64 `json:"wind_speed"`
		WindSpeedOfGust     float64 `json:"wind_speed_of_gust"`
		PrecipitationAmount float64 `json:"precipitation_amount"`
	} `json:"instant"`
	Predicted struct {
		Summary                    string  `json:"summary"`
		Confidence                 string  `json:"confidence"`
		AirTemperatureMax          float64 `json:"air_temperature_max"`
		AirTemperatureMin          float64 `json:"air_temperature_min"`
		PrecipitationAmount        float64 `json:"precipitation_amount"`
		PrecipitationAmountMax     float64 `json:"precipitation_amount_max"`
		PrecipitationAmountMin     float64 `json:"precipitation_amount_min"`
		ProbabilityOfPrecipitation float64 `json:"probability_of_precipitation"`
	} `json:"predicted"`
}

// Yr structure stores weather data for the next 10 days.
//
// Functionality: get, req
type Yr struct {
	Properties struct {
		Timeseries []struct {
			Time string `json:"time"`
			Data Data   `json:"data"`
		} `json:"timeseries"`
	} `json:"properties"`
}

// Data structure stores all weather data for a day.
type Data struct {
	Instant struct {
		Details struct {
			AirTemperature           float64 `json:"air_temperature"`
			CloudAreaFraction        float64 `json:"cloud_area_fraction"`
			DewPointTemperature      float64 `json:"dew_point_temperature"`
			RelativeHumidity         float64 `json:"relative_humidity"`
			WindFromDirection        float64 `json:"wind_from_direction"`
			WindSpeed                float64 `json:"wind_speed"`
			WindSpeedOfGust          float64 `json:"wind_speed_of_gust"`
		} `json:"details"`
	} `json:"instant"`
	Next12Hours struct {
		Summary struct {
			SymbolCode       string `json:"symbol_code"`
			SymbolConfidence string `json:"symbol_confidence"`
		} `json:"summary"`
		Details struct {
			ProbabilityOfPrecipitation float64 `json:"probability_of_precipitation"`
		} `json:"details"`
	} `json:"next_12_hours"`
	Next1Hours struct {
		Details struct {
			PrecipitationAmount        float64 `json:"precipitation_amount"`
		} `json:"details"`
	} `json:"next_1_hours"`
	Next6Hours struct {
		Details struct {
			AirTemperatureMax          float64 `json:"air_temperature_max"`
			AirTemperatureMin          float64 `json:"air_temperature_min"`
			PrecipitationAmount        float64 `json:"precipitation_amount"`
			PrecipitationAmountMax     float64 `json:"precipitation_amount_max"`
			PrecipitationAmountMin     float64 `json:"precipitation_amount_min"`
		} `json:"details"`
	} `json:"next_6_hours"`
}
