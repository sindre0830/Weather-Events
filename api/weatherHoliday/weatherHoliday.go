package weatherHoliday

import (
	"encoding/json"
	geocoords "main/api/geoCoords"
	"main/debug"
	"net/http"
)

type WeatherHoliday struct {
	Holiday string `json:"holiday"`
	Location string `json:"location"`
	Frequency int `json:"frequency"`
}

// Handler for the weather holiday webhook endpoint
func (weatherHoliday *WeatherHoliday) Handler(w http.ResponseWriter, r *http.Request) {
	// Decode body into struct
	err := json.NewDecoder(r.Body).Decode(&weatherHoliday)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"WeatherHoliday.Handler() -> Decoding body to struct",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	// Get the geocoords of the location
	var locationCoords geocoords.LocationCoords
	status, err := locationCoords.Handler(weatherHoliday.Location)
	if err != nil {
		debug.ErrorMessage.Update(
			status,
			"WeatherHoliday.Handler() -> LocationCoords.Handler() -> Getting main location info",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
}



