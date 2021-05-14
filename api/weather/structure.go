package weather

import "main/api/weatherData"

// Weather structure stores current and predicted weather data for a day and information about location.
//
// Functionality: Handler, get
type Weather struct {
	Longitude float64                        `json:"longitude"`
	Latitude  float64                        `json:"latitude"`
	Location  string                         `json:"location"`
	Updated   string                         `json:"updated"`
	Date	  string                         `json:"date"`
	Data 	  weatherData.WeatherDataForADay `json:"data"`
}
