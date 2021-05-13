package weatherEvent

import (
	"encoding/json"
	"main/api/holidaysData"
	"main/api/notification"
	"main/api/weather"
	"main/db"
	"main/debug"
	"main/dict"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"time"
)

// get handles a get request from the client.
func (weatherEvent *WeatherEvent) get(w http.ResponseWriter, r *http.Request) {
	//split URL path by '/' and branch if there aren't enough elements
	arrPath := strings.Split(r.URL.Path, "/")
	if len(arrPath) != 6 {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"WeatherEvent.get() -> Checking length of URL",
			"URL validation: either too many or too few arguments in URL path",
			"URL format. Expected format: '.../id'. Example: '.../1ab24db3",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//set id and check if it's specified by client
	id := arrPath[5]
	if id != "" {
		data, exist := db.DB.Get("weatherEvent", id)
		if !exist {
			debug.ErrorMessage.Update(
				http.StatusBadRequest,
				"WeatherEvent.get() -> Database.Get() -> finding document based on ID",
				"getting webhook: can't find id",
				"ID doesn't exist. Expected format: '.../id'. Example: '.../1ab24db3",
			)
			debug.ErrorMessage.Print(w)
			return
		}
		weatherEvent.readData(data["Container"].(interface{}))
		//update header to JSON and set HTTP code
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		//send output to user and branch if an error occured
		err := json.NewEncoder(w).Encode(&weatherEvent)
		if err != nil {
			debug.ErrorMessage.Update(
				http.StatusInternalServerError,
				"WeatherEvent.get() -> Sending data to user",
				err.Error(),
				"Unknown",
			)
			debug.ErrorMessage.Print(w)
			return
		}
	} else {
		arrData, err := db.DB.GetAll("weatherEvent")
		if err != nil {
			debug.ErrorMessage.Update(
				http.StatusInternalServerError,
				"WeatherEvent.get() -> Database.GetAll() -> Getting all documents",
				err.Error(),
				"Unknown",
			)
			debug.ErrorMessage.Print(w)
			return
		}
		var arrWeatherEvent []WeatherEvent
		for _, rawData := range arrData {
			data := rawData["Container"].(interface{})
			weatherEvent.readData(data)
			arrWeatherEvent = append(arrWeatherEvent, *weatherEvent)
		}
		//update header to JSON and set HTTP code
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		//send output to user and branch if an error occured
		err = json.NewEncoder(w).Encode(&arrWeatherEvent)
		if err != nil {
			debug.ErrorMessage.Update(
				http.StatusInternalServerError,
				"WeatherEvent.get() -> Sending data to user",
				err.Error(),
				"Unknown",
			)
			debug.ErrorMessage.Print(w)
			return
		}
	}
}

// post handles a post request from the client.
func (weatherEvent *WeatherEvent) post(w http.ResponseWriter, r *http.Request) {
	//read input from client and branch if an error occurred
	var weatherEventInput WeatherEventInput
	err := json.NewDecoder(r.Body).Decode(&weatherEventInput)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"WeatherEvent.post() -> Parsing data from client",
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
			"WeatherEvent.post() -> Checking if URL is valid",
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
			"WeatherEvent.post() -> Checking if URL is valid",
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
			"WeatherEvent.post() -> Checking if timeout value is valid",
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
			"WeatherEvent.post() -> Checking if trigger value is valid",
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
			"WeatherEvent.post() -> Checking if location is valid",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	recorder := httptest.NewRecorder()
	weather.Handler(recorder, req)
	if recorder.Code != http.StatusOK {
		debug.ErrorMessage.Update(
			http.StatusNotFound,
			"WeatherEvent.post() -> Checking if location is valid",
			err.Error(),
			"Location not found. Example: 'Oslo'",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//convert holiday to date if it is inputted
	weatherEventInput.checkIfHoliday(w)
	//check if date is valid
	if !weatherEventInput.checkDate() {
		debug.ErrorMessage.Update(
			http.StatusNotFound,
			"WeatherEvent.post() -> WeatherEvent.checkDate() -> Checking if date is valid",
			"invalid date: date is either wrong format or not within scope",
			"Check that the format of the date is YYYY-MM-DD and that it is within timeframe.",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//set data
	weatherEvent.Date = weatherEventInput.Date
	weatherEvent.Location = weatherEventInput.Location
	weatherEvent.URL = weatherEventInput.URL
	weatherEvent.Frequency = weatherEventInput.Frequency
	weatherEvent.Timeout = weatherEventInput.Timeout
	//send data to database
	var data db.Data
	data.Container = weatherEvent
	_, id, err := db.DB.Add("weatherEvent", "", data)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"WeatherEvent.post() -> Database.Add() -> Adding data to database",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	weatherEvent.ID = id
	//create feedback message to send to client and branch if an error occurred
	var feedback notification.Feedback
	feedback.Update(
		http.StatusCreated,
		"Webhook successfully created for '" + weatherEvent.URL + "'",
		weatherEvent.ID,
	)
	err = feedback.Print(w)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"WeatherEvent.post() -> Feedback.print() -> Sending feedback to client",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//start loop
	go weatherEvent.callHook()
}

// delete handles a delete request from the client.
func (weatherEvent *WeatherEvent) delete(w http.ResponseWriter, r *http.Request) {
	//split URL path by '/' and branch if there aren't enough elements
	arrPath := strings.Split(r.URL.Path, "/")
	if len(arrPath) != 6 {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"WeatherEvent.delete() -> Checking length of URL",
			"URL validation: either too many or too few arguments in URL path",
			"URL format. Expected format: '.../id'. Example: '.../1ab24db3",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//set id and check if it's specified by client
	id := arrPath[5]
	err := db.DB.Delete("weatherEvent", id)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"WeatherEvent.delete() -> Database.Delete() -> Deleting document based on ID",
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
			"WeatherEvent.delete() -> Feedback.print() -> Sending feedback to client",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
}

// readData parses data from database to WeatherEvent structure format.
func (weatherEvent *WeatherEvent) readData(data interface{}) {
	rawData := data.(map[string]interface{})
	weatherEvent.ID = rawData["ID"].(string)
	weatherEvent.Date = rawData["Date"].(string)
	weatherEvent.Location = rawData["Location"].(string)
	weatherEvent.URL = rawData["URL"].(string)
	weatherEvent.Frequency = rawData["Frequency"].(string)
	weatherEvent.Timeout = rawData["Timeout"].(int64)
}

// checkDate checks if date is valid and within timeframe.
func (weatherEventInput *WeatherEventInput) checkDate() bool {
	date, err := time.Parse("2006-01-02", weatherEventInput.Date)
	if err != nil {
		return false
	}
	if weatherEventInput.Date == time.Now().Format("2006-01-02") {
		return true
	} else {
		return time.Now().Before(date)
	}
}

// checkIfHoliday converts holiday to date.
func (weatherEventInput *WeatherEventInput) checkIfHoliday(w http.ResponseWriter) {
	// Parse date to see if it is a date or a holiday
	_, err := time.Parse("2006-01-02", weatherEventInput.Date)
	if err != nil {
		// It is a holiday, replace holiday with date
		// Get a map of all the country's holidays
		var holidaysMap = make(map[string]interface{})
		holidaysMap, status, err := holidaysData.Handler(weatherEventInput.Location)
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
		weatherEventInput.Date = strings.Title(strings.ToLower(weatherEventInput.Date))

		// Check if the holiday exists in the selected country
		date, ok := holidaysMap[weatherEventInput.Date]
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

		weatherEventInput.Date = date.(string)
	}
}
