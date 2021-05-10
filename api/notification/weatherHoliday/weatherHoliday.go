package weatherHoliday

import (
	"encoding/json"
	"main/api/holidaysData"
	"main/api/notification"
	"main/db"
	"main/debug"
	"net/http"
	"net/url"
	"strings"
)

// WeatherHolidayInput structure, stores information from the user about the webhook
type WeatherHolidayInput struct {
	Holiday 	string `json:"holiday"`
	Location 	string `json:"location"`
	URL 		string `json:"url"`
	Frequency 	string `json:"frequency"`
	Timeout 	int64 `json:"timeout"`
}

// WeatherHoliday structure, stores information about the webhook added to the database
type WeatherHoliday struct {
	ID        	string `json:"id"`
	Date 		string `json:"date"`
	Holiday 	string `json:"holiday"`
	Location 	string `json:"location"`
	URL 		string `json:"url"`
	Frequency 	string `json:"frequency"`
	Timeout 	int64 `json:"timeout"`
}


// POST handles a POST request from the http request.
func (weatherHoliday *WeatherHoliday) POST(w http.ResponseWriter, r *http.Request) {
	var weatherHolidayInput WeatherHolidayInput

	// Decode body into weatherHoliday struct
	err := json.NewDecoder(r.Body).Decode(&weatherHoliday)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"WeatherHoliday.POST() -> Decoding body to struct",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	// Check if the URL the user sent is valid
	parsedURL, err := url.ParseRequestURI(weatherHolidayInput.URL)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"WeatherHoliday.POST() -> Checking if URL is valid",
			err.Error(),
			"Not valid URL. Example 'http://google.com/'",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	// Branch if the schema in the URL is incorrect
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"WeatherHoliday.POST() -> Checking if URL is valid",
			"url validation: schema is incorrect, should be 'http' or 'https'",
			"Not valid URL. Example 'http://google.com/'",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	// Check if the timeout sent by the user is valid
	if weatherHolidayInput.Timeout < 15 || weatherHolidayInput.Timeout > 86400 {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"WeatherHoliday.POST() -> Checking if timeout value is valid",
			"timeout validation: value isn't within scope",
			"Timeout value has to be larger then 15 and less then 86400(24 hours) seconds.",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	// Check if the trigger sent by the user is valid
	weatherHolidayInput.Frequency = strings.ToUpper(weatherHolidayInput.Frequency)
	if weatherHolidayInput.Frequency != "EVERY_DAY" && weatherHolidayInput.Frequency != "ON_DATE" {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"WeatherHoliday.POST() -> Checking if trigger value is valid",
			"trigger validation: trigger is not 'EVERY_DAY' or 'ON_DATE'",
			"Not valid trigger. Example 'ON_DATE'",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	// Get a map of all the country's holidays
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

	// Make the first letter of each word uppercase to match the format in holidaysMap
	weatherHoliday.Holiday = strings.Title(strings.ToLower(weatherHoliday.Holiday))

	// Check if the holiday exists in the selected country
	date, ok := holidaysMap[weatherHoliday.Holiday]
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

	// Set data to database struct
	weatherHoliday.Date = date.(string)
	weatherHoliday.Holiday = weatherHolidayInput.Holiday
	weatherHoliday.Location = weatherHolidayInput.Location
	weatherHoliday.URL = weatherHolidayInput.URL
	weatherHoliday.Frequency = weatherHolidayInput.Frequency
	weatherHoliday.Timeout = weatherHolidayInput.Timeout

	// Add data to database
	var data db.Data
	data.Container = weatherHoliday

	_, id, err := db.DB.Add("weatherHoliday", "", data)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"WeatherHoliday.POST() -> Database.Add() -> Adding data to database",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	// Create feedback message and print it to the user
	var feedback notification.Feedback
	feedback.Update(
		http.StatusCreated,
		"Webhook successfully created for '" + weatherHoliday.URL + "'",
		id,
	)

	err = feedback.Print(w)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"WeatherEvent.POST() -> Feedback.print() -> Sending feedback to client",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
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


