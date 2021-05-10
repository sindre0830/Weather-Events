package weatherEvent

import (
	"encoding/json"
	"main/api/notification"
	"main/api/weather"
	"main/db"
	"main/debug"
	"main/dict"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

type WeatherEventInput struct {
	Date      string `json:"date"`
	Location  string `json:"location"`
	URL       string `json:"url"`
	Frequency string `json:"frequency"`
	Timeout   int    `json:"timeout"`
}

type WeatherEvent struct {
	ID        string `json:"id"`
	Date      string `json:"date"`
	Location  string `json:"location"`
	URL       string `json:"url"`
	Frequency string `json:"frequency"`
	Timeout   int    `json:"timeout"`
}

// POST handles a POST request from the http request.
func (weatherEvent *WeatherEvent) POST(w http.ResponseWriter, r *http.Request) {
	//read input from client and branch if an error occurred
	var weatherEventInput WeatherEventInput
	err := json.NewDecoder(r.Body).Decode(&weatherEventInput)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusBadRequest, 
			"WeatherEvent.POST() -> Parsing data from client",
			err.Error(),
			"Wrong JSON format sent.",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//check if URL is valid (very simple check) and branch if an error occurred
	parsedURL, err := url.ParseRequestURI(weatherEventInput.URL)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"WeatherEvent.POST() -> Checking if URL is valid",
			err.Error(),
			"Not valid URL. Example 'http://google.com/'",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//branch if the schema in the URL is incorrect
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"WeatherEvent.POST() -> Checking if URL is valid",
			"url validation: schema is incorrect, should be 'http' or 'https'",
			"Not valid URL. Example 'http://google.com/'",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//check if timeout is valid and return an error if it isn't
	if weatherEventInput.Timeout < 15 || weatherEventInput.Timeout > 86400 {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"WeatherEvent.POST() -> Checking if timeout value is valid",
			"timeout validation: value isn't within scope",
			"Timeout value has to be larger then 15 and less then 86400(24 hours) seconds.",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//check if trigger is valid and return an error if it isn't
	weatherEventInput.Frequency = strings.ToUpper(weatherEventInput.Frequency)
	if weatherEventInput.Frequency != "EVERY_DAY" && weatherEventInput.Frequency != "ON_DATE" {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"WeatherEvent.POST() -> Checking if trigger value is valid",
			"trigger validation: trigger is not 'EVERY_DAY' or 'ON_DATE'",
			"Not valid trigger. Example 'ON_DATE'",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//validate parameters and branch if an error occurred
	var weather weather.Weather
	req, err := http.NewRequest("GET", dict.GetWeatherURL(weatherEventInput.Location, ""), nil)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"WeatherEvent.POST() -> Checking if location is valid",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
	}
	recorder := httptest.NewRecorder()
	weather.Handler(recorder, req)
	if recorder.Code != http.StatusOK {
		debug.ErrorMessage.Update(
			http.StatusNotFound,
			"WeatherEvent.POST() -> Checking if location is valid",
			err.Error(),
			"Location not found. Example: 'Oslo'",
		)
		debug.ErrorMessage.Print(w)
	}
	//set data
	weatherEvent.Date = weatherEventInput.Date
	weatherEvent.Location = weatherEventInput.Location
	weatherEvent.URL = weatherEventInput.URL
	weatherEvent.Frequency = weatherEventInput.Frequency
	weatherEvent.Timeout = weatherEventInput.Timeout
	var data db.Data
	data.Container = weatherEvent
	_, id, err := db.DB.Add("weatherEvent", "", data)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"WeatherEvent.POST() -> Database.Add() -> Adding data to database",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//create feedback message to send to client and branch if an error occurred
	var feedback notification.Feedback
	feedback.Update(
		http.StatusCreated, 
		"Webhook successfully created for '" + weatherEvent.URL + "'",
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
