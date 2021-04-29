package weather

type Weather struct {
	Location  string  `json:"location"`
	Longitude float64 `json:"longitude"`
	Latiude   float64 `json:"latiude"`
	Country   string  `json:"country"`
	Updated   string  `json:"updated"`
	Data      struct {
		Now struct {
			AirTemperature      float64 `json:"air_temperature"`
			CloudAreaFraction   float64 `json:"cloud_area_fraction"`
			DewPointTemperature float64 `json:"dew_point_temperature"`
			RelativeHumidity    float64 `json:"relative_humidity"`
			WindFromDirection   float64 `json:"wind_from_direction"`
			WindSpeed           float64 `json:"wind_speed"`
			WindSpeedOfGust     float64 `json:"wind_speed_of_gust"`
			PrecipitationAmount float64 `json:"precipitation_amount"`
		} `json:"now"`
		Today struct {
			Summary                    string  `json:"summary"`
			Confidence                 string  `json:"confidence"`
			AirTemperatureMax          float64 `json:"air_temperature_max"`
			AirTemperatureMin          float64 `json:"air_temperature_min"`
			PrecipitationAmount        float64 `json:"precipitation_amount"`
			PrecipitationAmountMax     float64 `json:"precipitation_amount_max"`
			PrecipitationAmountMin     float64 `json:"precipitation_amount_min"`
			ProbabilityOfPrecipitation float64 `json:"probability_of_precipitation"`
		} `json:"today"`
	} `json:"data"`
}
