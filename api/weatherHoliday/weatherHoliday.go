package weatherHoliday

import (
	"encoding/json"
	"main/api/countryData"
	"main/api/geocoords"
	"main/api/holidaysData"
	"main/db"
	"main/debug"
	"net/http"
	"strings"
)

// Request
type WeatherHoliday struct {
	Holiday string `json:"holiday"`
	Location string `json:"location"`
	URL string `json:"url"`
	Frequency string `json:"frequency"`		// Every day or on date
	Timeout int `json:"timeout"`			// Hours
	ID string `json:"id"`
}

// Register a webhook
func (weatherHoliday *WeatherHoliday) Register(w http.ResponseWriter, r *http.Request) {
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
			"WeatherHoliday.Register() -> LocationCoords.Handler() -> Getting location info",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	// Get country and format it correctly
	address := strings.Split(locationCoords.Address, ", ")
	country := address[len(address)-1]

	// Get country code
	var countryInfo countryData.Information

	status, err, countryCode := countryInfo.Handler(country)
	if err != nil {
		debug.ErrorMessage.Update(
			status,
			"WeatherHoliday.Register() -> CountryData.handler() -> Getting country code",
			err.Error(),
			"Selected country is not valid",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	// Get the country's holidays
	var holidaysMap = make(map[string]interface{})
	holidaysMap, status, err = holidaysData.Handler(countryCode)
	if err != nil {
		debug.ErrorMessage.Update(
			status,
			"WeatherHoliday.Register() -> holidaysData.Handler() - > Getting information about the country's holidays",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	// Make the first letter of each word uppercase
	weatherHoliday.Holiday = strings.Title(weatherHoliday.Holiday)

	// Check if the holiday exists in the selected country
	_, ok := holidaysMap[weatherHoliday.Holiday]
	if !ok {
		http.Error(w, "The selected holiday is not valid", http.StatusBadRequest)
		return
	}

	// Check if the frequency field is valid
	if weatherHoliday.Frequency != "ON_DATE" && weatherHoliday.Frequency != "EVERY_DAY" {
		http.Error(w, "The selected frequency is not valid. Try writing either 'ON_DATE' or 'EVERY_DAY'", http.StatusBadRequest)
		return
	}

	// Add webhook to the database
	var dataDB db.Data
	dataDB.Container = weatherHoliday

	_, id, err := db.DB.Add("WeatherHoliday", "", dataDB)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"WeatherHoliday.Register() -> db.Add() -> Adding webhook to the database",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	weatherHoliday.ID = id

	http.Error(w, "Webhook registered", http.StatusCreated)
}

// Delete a webhook
func (weatherHoliday *WeatherHoliday) Delete(w http.ResponseWriter, r *http.Request) {
	// Parse URL path and ensure that the formatting is correct
	path := strings.Split(r.URL.Path, "/")
	if len(path) != 6 {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"WeatherHoliday.Delete() -> Parsing URL",
			"url validation: either too many or too few arguments in url path",
			"URL format. Remember to add an ID at the end of the path",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	// Get webhook ID
	id := path[len(path)-1]

	err := db.DB.Delete(id)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"WeatherHoliday.Delete() -> db.Delete() -> Deleting webhook from the database",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	http.Error(w, "Webhook successfully deleted", http.StatusNoContent)
}


