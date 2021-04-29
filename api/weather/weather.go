package weather

import (
	"fmt"
	"main/api/weatherData"
	"net/http"
)

// Weather structure stores current and predicted weather data for a day and information about location.
//
// Functionality: get
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

// get will get data for structure.
func (weather *Weather) get(lat float64, lon float64) (int, error) {
	//convert coordinates to string
	strLat := fmt.Sprintf("%f", lat)
	strLon := fmt.Sprintf("%f", lon)
	//get weather data and branch if an error occurred
	var weatherData weatherData.WeatherData
	status, err := weatherData.Handler(strLat, strLon)
	if err != nil {
		return status, err
	}
	//set data in structure
	weather.Updated = weatherData.Updated
	weather.Data.Now = weatherData.Now
	weather.Data.Today = weatherData.Today
	return http.StatusOK, nil
}
