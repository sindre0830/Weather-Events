package weather

import (
	"encoding/json"
	"main/api"
	"main/api/diag"
	"main/api/notification"
	"main/api/weatherDetails"
	"main/debug"
	"main/dict"
	"main/storage"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

// get handles a get request from the client.
func (weather *Weather) get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//split URL path by '/' and branch if there aren't enough elements
	arrPath := strings.Split(r.URL.Path, "/")
	if len(arrPath) != 6 {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"Weather.get() -> Checking length of URL",
			"URL validation: either too many or too few arguments in URL path",
			"URL format. Expected format: '.../id'. Example: '.../1ab24db3",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//set id and check if it's specified by client
	id := arrPath[5]
	if id != "" {
		data, exist := storage.Firebase.Get(dict.WEATHER_COLLECTION, id)
		if !exist {
			debug.ErrorMessage.Update(
				http.StatusBadRequest,
				"Weather.get() -> Database.Get() -> finding document based on ID",
				"getting webhook: can't find id",
				"ID doesn't exist. Expected format: '.../id'. Example: '.../1ab24db3",
			)
			debug.ErrorMessage.Print(w)
			return
		}
		weather.readData(data["Container"].(interface{}))
		//send data to client and branch if an error occured
		err := api.SendData(w, weather, http.StatusOK)
		if err != nil {
			debug.ErrorMessage.Update(
				http.StatusInternalServerError,
				"Weather.get() -> Sending data to user",
				err.Error(),
				"Unknown",
			)
			debug.ErrorMessage.Print(w)
			return
		}
	} else {
		arrData, err := storage.Firebase.GetAll(dict.WEATHER_COLLECTION)
		if err != nil {
			debug.ErrorMessage.Update(
				http.StatusInternalServerError,
				"Weather.get() -> Database.GetAll() -> Getting all documents",
				err.Error(),
				"Unknown",
			)
			debug.ErrorMessage.Print(w)
			return
		}
		var arrWeather []Weather
		for _, rawData := range arrData {
			data := rawData["Container"].(interface{})
			weather.readData(data)
			arrWeather = append(arrWeather, *weather)
		}
		//send data to client and branch if an error occured
		err = api.SendData(w, arrWeather, http.StatusOK)
		if err != nil {
			debug.ErrorMessage.Update(
				http.StatusInternalServerError,
				"Weather.get() -> Sending data to user",
				err.Error(),
				"Unknown",
			)
			debug.ErrorMessage.Print(w)
			return
		}
	}
}

// post handles a post request from the client.
func (weather *Weather) post(w http.ResponseWriter, r *http.Request) {
	// turn json object into struct
	var weatherInput WeatherInput
	err := json.NewDecoder(r.Body).Decode(&weatherInput)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"Weather.Handler() -> Decoding body",
			err.Error(),
			"Improper formatting of request body.",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//set data
	weather.Location = weatherInput.Location
	weather.URL = weatherInput.URL
	weather.Timeout = weatherInput.Timeout
	//check if URL is valid (very simple check) and branch if an error occurred
	parsedURL, err := url.ParseRequestURI(weather.URL)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"Weather.post() -> Checking if URL is valid",
			err.Error(),
			"Not valid URL. Example 'http://google.com/'",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	// branch if the schema in the URL is incorrect
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"Weather.post() -> Checking if URL is valid",
			"url validation: schema is incorrect, should be 'http' or 'https'",
			"Not valid URL. Example 'http://google.com/'",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//check if timeout is valid and return an error if it isn't
	if weather.Timeout < 15 || weather.Timeout > 86400 {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"Weather.post() -> Checking if timeout value is valid",
			"timeout validation: value isn't within scope",
			"Timeout value has to be larger then 15 and less then 86400(24 hours) seconds.",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//validate parameters and branch if an error occurred
	var weatherDetails weatherDetails.WeatherDetails
	req, err := http.NewRequest("GET", dict.GetWeatherDetailsURL(weather.Location, ""), nil)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"Weather.post() -> Checking if location is valid",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	recorder := httptest.NewRecorder()
	weatherDetails.Handler(recorder, req)
	if recorder.Code != http.StatusOK {
		debug.ErrorMessage.Update(
			http.StatusNotFound,
			"Weather.post() -> Checking if location is valid",
			"validating location: couldn't find location",
			"Location not found. Example: 'Oslo'",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//send data to database
	var data storage.Data
	data.Container = weather
	_, id, err := storage.Firebase.Add(dict.WEATHER_COLLECTION, "", data)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"Weather.post() -> Database.Add() -> Adding data to database",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	weather.ID = id
	//create feedback message to send to client and branch if an error occurred
	var feedback notification.Feedback
	feedback.Update(
		http.StatusCreated,
		"Webhook successfully created for '" + weather.URL + "'",
		weather.ID,
	)
	err = feedback.Print(w)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"Weather.post() -> Feedback.print() -> Sending feedback to client",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//start loop
	go weather.callHook()
	//add hook amount to diag
	diag.HookAmount++
}

// delete handles a delete request from the client.
func (weather *Weather) delete(w http.ResponseWriter, r *http.Request) {
	//split URL path by '/' and branch if there aren't enough elements
	arrPath := strings.Split(r.URL.Path, "/")
	if len(arrPath) != 6 {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"Weather.delete() -> Checking length of URL",
			"URL validation: either too many or too few arguments in URL path",
			"URL format. Expected format: '.../id'. Example: '.../1ab24db3",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//set id and check if it's specified by client
	id := arrPath[5]
	err := storage.Firebase.Delete(dict.WEATHER_COLLECTION, id)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"Weather.delete() -> Database.Delete() -> Deleting document based on ID",
			err.Error(),
			"ID doesn't exist. Expected format: '.../id'. Example: '.../1ab24Db3",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//create feedback message to send to client and branch if an error occurred
	var feedback notification.Feedback
	feedback.Update(
		http.StatusOK,
		"Webhook successfully deleted",
		id,
	)
	err = feedback.Print(w)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"Weather.delete() -> Feedback.print() -> Sending feedback to client",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
}

// readData parses data from database to Weather structure format.
func (weather *Weather) readData(data interface{}) {
	rawData := data.(map[string]interface{})
	weather.ID = rawData["ID"].(string)
	weather.Location = rawData["Location"].(string)
	weather.URL = rawData["URL"].(string)
	weather.Timeout = rawData["Timeout"].(int64)
}