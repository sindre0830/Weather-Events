package weatherDetails

import "main/api/weatherData"

// WeatherDetails structure stores current and predicted weather data for a day and information about location.
//
// Functionality: Handler, get
type WeatherDetails struct {
	Longitude float64                        `json:"longitude"`
	Latitude  float64                        `json:"latitude"`
	Location  string                         `json:"location"`
	Updated   string                         `json:"updated"`
	Date	  string                         `json:"date"`
	Data 	  weatherData.WeatherDataForADay `json:"data"`
}
