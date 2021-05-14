package weatherEvent

import (
	"encoding/json"
	"main/api/diag"
	"main/api/eventData"
	"main/api/holidaysData"
	"main/api/notification"
	"main/api/weatherDetails"
	"main/debug"
	"main/dict"
	"main/storage"
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
		data, exist := storage.Firebase.Get("weatherEvent", id)
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
		arrData, err := storage.Firebase.GetAll("weatherEvent")
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

func (weatherEvent *WeatherEvent) postHandler(w http.ResponseWriter, r *http.Request) {
	//parse url and branch if an error occurred
	arrPath := strings.Split(r.URL.Path, "/")
	if len(arrPath) != 6 {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"WeatherEvent.postHandler() -> Parsing URL",
			"url validation: either too many or too few arguments in url path",
			"URL format.",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	event := strings.ToLower(arrPath[5])
	switch event {
		case "":
			//read input from client and branch if an error occurred
			var weatherEventDefault WeatherEventDefault
			err := json.NewDecoder(r.Body).Decode(&weatherEventDefault)
			if err != nil {
				debug.ErrorMessage.Update(
					http.StatusBadRequest,
					"WeatherEvent.postHandler() -> Parsing data from client",
					err.Error(),
					"Wrong JSON format sent. See README for template and example.",
				)
				debug.ErrorMessage.Print(w)
				return
			}
			//set data
			weatherEvent.Date = weatherEventDefault.Date
			weatherEvent.Location = weatherEventDefault.Location
			weatherEvent.URL = weatherEventDefault.URL
			weatherEvent.Frequency = weatherEventDefault.Frequency
			weatherEvent.Timeout = weatherEventDefault.Timeout
		case "holiday":
			//read input from client and branch if an error occurred
			var weatherEventHoliday WeatherEventHoliday
			err := json.NewDecoder(r.Body).Decode(&weatherEventHoliday)
			if err != nil {
				debug.ErrorMessage.Update(
					http.StatusBadRequest,
					"WeatherEvent.postHandler() -> Parsing data from client",
					err.Error(),
					"Wrong JSON format sent. See README for template and example.",
				)
				debug.ErrorMessage.Print(w)
				return
			}
			//get a map of all the country's holidays
			var holidaysMap = make(map[string]interface{})
			holidaysMap, status, err := holidaysData.Handler(weatherEventHoliday.Location)
			if err != nil {
				debug.ErrorMessage.Update(
					status,
					"WeatherEvent.postHandler() -> WeatherHoliday.Register() -> holidaysData.Handler() - > Getting information about the country's holidays",
					err.Error(),
					"Unknown",
				)
				debug.ErrorMessage.Print(w)
				return
			}
			//make the first letter of each word uppercase to match the format in holidaysMap
			weatherEventHoliday.Holiday = strings.Title(strings.ToLower(weatherEventHoliday.Holiday))
			//check if the holiday exists in the selected country
			if date, ok := holidaysMap[weatherEventHoliday.Holiday].(string); !ok {
				debug.ErrorMessage.Update(
					http.StatusBadRequest,
					"WeatherHoliday.Register() -> Checking if a holiday exists in a country",
					"invalid holiday: the holiday is not valid in the selected country",
					"Not a real holiday. Check your spelling and make sure it is the english name.",
				)
				debug.ErrorMessage.Print(w)
				return
			} else {
				weatherEvent.Date = date
			}
			//set data
			weatherEvent.Location = weatherEventHoliday.Location
			weatherEvent.URL = weatherEventHoliday.URL
			weatherEvent.Frequency = weatherEventHoliday.Frequency
			weatherEvent.Timeout = weatherEventHoliday.Timeout
		case "ticket":
			//read input from client and branch if an error occurred
			var weatherEventTicketmaster WeatherEventTicketmaster
			err := json.NewDecoder(r.Body).Decode(&weatherEventTicketmaster)
			if err != nil {
				debug.ErrorMessage.Update(
					http.StatusBadRequest,
					"WeatherEvent.postHandler() -> Parsing data from client",
					err.Error(),
					"Wrong JSON format sent. See README for template and example.",
				)
				debug.ErrorMessage.Print(w)
				return
			}
			var ticketMaster eventData.FirebaseStore
			status, err := ticketMaster.Handler(weatherEventTicketmaster.Ticket)
			if err != nil {
				debug.ErrorMessage.Update(
					status,
					"WeatherEvent.postHandler() -> FirebaseStore.Handler -> Getting ticket information",
					err.Error(),
					"Invalid ticket.",
				)
				debug.ErrorMessage.Print(w)
				return
			}
			//set data
			weatherEvent.Date = ticketMaster.Localdate
			weatherEvent.Location = ticketMaster.Name
			weatherEvent.URL = weatherEventTicketmaster.URL
			weatherEvent.Frequency = weatherEventTicketmaster.Frequency
			weatherEvent.Timeout = weatherEventTicketmaster.Timeout
		default:
			debug.ErrorMessage.Update(
				http.StatusBadRequest, 
				"WeatherEvent.postHandler() -> Validating event type",
				"event validation: event doesn't exist",
				"Event not implemented. List of possible events: 'Holiday' where date is defined for client, 'Ticketmaster' where date and location is defined for client, '' where client has to define location and date.",
			)
			debug.ErrorMessage.Print(w)
			return
	}
	weatherEvent.post(w, r)
}

// post handles a post request from the client.
func (weatherEvent *WeatherEvent) post(w http.ResponseWriter, r *http.Request) {
	//check if URL is valid (very simple check) and branch if an error occurred
	parsedURL, err := url.ParseRequestURI(weatherEvent.URL)
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
	if weatherEvent.Timeout < 15 || weatherEvent.Timeout > 86400 {
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
	if weatherEvent.Frequency != "EVERY_DAY" && weatherEvent.Frequency != "ON_DATE" {
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
	var weatherDetails weatherDetails.WeatherDetails
	req, err := http.NewRequest("GET", dict.GetWeatherURL(weatherEvent.Location, ""), nil)
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
	weatherDetails.Handler(recorder, req)
	if recorder.Code != http.StatusOK {
		debug.ErrorMessage.Update(
			http.StatusNotFound,
			"WeatherEvent.post() -> Checking if location is valid",
			"validating location: couldn't find location",
			"Location not found. Example: 'Oslo'",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//check if date is valid
	if !weatherEvent.checkDate() {
		debug.ErrorMessage.Update(
			http.StatusNotFound,
			"WeatherEvent.post() -> WeatherEvent.checkDate() -> Checking if date is valid",
			"invalid date: date is either wrong format or not within scope",
			"Check that the format of the date is YYYY-MM-DD and that it is within timeframe.",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//send data to database
	var data storage.Data
	data.Container = weatherEvent
	_, id, err := storage.Firebase.Add("weatherEvent", "", data)
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
	//add hook amount to diag
	diag.HookAmount++
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
	err := storage.Firebase.Delete("weatherEvent", id)
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
func (weatherEvent *WeatherEvent) checkDate() bool {
	date, err := time.Parse("2006-01-02", weatherEvent.Date)
	if err != nil {
		return false
	}
	if weatherEvent.Date == time.Now().Format("2006-01-02") {
		return true
	} else {
		return time.Now().Before(date)
	}
}
