package weatherHoliday

import (
	"encoding/json"
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

	// Get the country's holidays
	var holidaysMap = make(map[string]interface{})
	holidaysMap, status, err := holidaysData.Handler(weatherHoliday.Location)
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
	weatherHoliday.Holiday = strings.Title(strings.ToLower(weatherHoliday.Holiday))

	// Check if the holiday exists in the selected country
	_, ok := holidaysMap[weatherHoliday.Holiday]
	if !ok {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"WeatherHoliday.Register() -> Checking if a holiday exists in a country",
			"invalid holiday: the holiday is not valid in the selected country",
			"Not a real holiday. Check your spelling and make sure it is the english name.",
		)
		debug.ErrorMessage.Print(w)
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

	http.Error(w, "Webhook registered with ID " + id, http.StatusCreated)
}

// Delete a webhook
func (weatherHoliday *WeatherHoliday) Delete(w http.ResponseWriter, r *http.Request) {
	// Parse URL path and ensure that the formatting is correct
	path := strings.Split(r.URL.Path, "/")
	if len(path) != 7 {
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

	err := db.DB.Delete("notifications", id)
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


