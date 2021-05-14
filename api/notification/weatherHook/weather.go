package weatherHook

import (
	"encoding/json"
	"main/api/diag"
	"main/api/notification"
	"main/api/weather"
	"main/debug"
	"main/dict"
	"main/storage"
	"net/http"
	"net/http/httptest"
	"net/url"
)

// name of the firestore database for this webhook
var hookdb = "weatherHookDB"

/**
* HandlerPost
* Handles POST method requests from client.
**/
func (weatherHook *WeatherHook) HandlerPost(w http.ResponseWriter, r *http.Request) {
	// turn json object into struct
	var weatherHookInput WeatherHookInput
	err := json.NewDecoder(r.Body).Decode(&weatherHookInput)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"WeatherHook.Handler() -> Decoding body",
			err.Error(),
			"Improper formatting of request body.",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	parsedURL, err := url.ParseRequestURI(weatherHookInput.URL)
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
	// branch if the schema in the URL is incorrect
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
	// validate parameters and branch if an error occurred
	var weather weather.Weather
	req, err := http.NewRequest("GET", dict.GetWeatherURL(weatherHookInput.Location, ""), nil)
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
	// check if timeout is valid and return an error if it isn't - timeout in hours
	if weatherHookInput.Timeout < 1 || weatherHookInput.Timeout > 72 {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"WeatherEvent.POST() -> Checking if timeout value is valid",
			"timeout validation: value isn't within scope",
			"Timeout value has to be at keast 1 and no more than 72 hours.",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	weatherHook.Location = weatherHookInput.Location
	weatherHook.Timeout = weatherHookInput.Timeout
	weatherHook.URL = weatherHookInput.URL
	// add it to firebase
	// how do we want to handle ID, getting/passing to user? Currently getting but not checking dupes.
	var data storage.Data
	data.Container = weatherHook
	_, id, err := storage.Firebase.Add(hookdb, "", data)
	// return ID if successful, otherwise error
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"WeatherHook -> MethodHandler() -> weatherHook.HandlerPost() -> Adding webhook to database.",
			"Database error: failed to add webhook!",
			"Improper formatting of webhook.",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	weatherHook.ID = id
	// start loop
	go weatherHook.callLoop()
	//add hook amount to diag
	diag.HookAmount++
	// create feedback message to send to client and branch if an error occurred
	var feedback notification.Feedback
	feedback.Update(
		http.StatusCreated,
		"Webhook successfully created for '"+weatherHook.URL+"'",
		weatherHook.ID,
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

/**
* HandlerGet
* Handles GET method requests from clients
**/
func (weatherHook *WeatherHook) HandlerGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// get query with ID
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusNotFound,
			"WeatherHook -> MethodHandler() -> weatherHook.HandlerGet() -> Parsing URL query.",
			"Parse error: Failed when handling query!",
			"",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	id := params["id"][0]

	// check for ID in firestore - DB function
	data, exist := storage.Firebase.Get(hookdb, id)
	// extract from data
	if exist {
		err := weatherHook.ReadData(data["Container"].(interface{}))
		if err != nil {
			debug.ErrorMessage.Update(
				http.StatusNotFound,
				"WeatherHook -> MethodHandler() -> weatherHook.HandlerGet() -> Reading data from firebase.",
				"Database error: Failed when reading from database!",
				"",
			)
			debug.ErrorMessage.Print(w)
			return
		}
		weatherHook.ID = id
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(&weatherHook)
		// send output to user and branch if an error occured
		if err != nil {
			debug.ErrorMessage.Update(
				http.StatusInternalServerError,
				"WeatherEvent.GET() -> Sending data to user",
				err.Error(),
				"Unknown",
			)
			debug.ErrorMessage.Print(w)
			return
		}
	} else {
		debug.ErrorMessage.Update(
			http.StatusNotFound,
			"WeatherHook -> MethodHandler() -> weatherHook.HandlerGet() -> Finding ID",
			"Database error: ID not found!",
			"Bad ID entered OR wrong method: GET.",
		)
		debug.ErrorMessage.Print(w)
	}
	return
}

/**
* HandlerDelete
* Handles DELETE method requests from client.
**/
func (weatherHook *WeatherHook) HandlerDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// get query with ID
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusNotFound,
			"WeatherHook -> MethodHandler() -> weatherHook.HandlerGet() -> Parsing URL query.",
			"Parse error: Failed when handling query!",
			"",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	id := params["id"][0]
	// delete
	err = storage.Firebase.Delete(hookdb, id)

	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"WeatherHook -> MethodHandler() -> weatherHook.HandlerDelete() -> Deleting webhook",
			"Database error on deleting webhook! Bad ID entered OR wrong method: DELETE.",
			"",
		)
		debug.ErrorMessage.Print(w)
		return
	} else {
		debug.ErrorMessage.Update(
			http.StatusOK,
			"WeatherHook -> MethodHandler() -> weatherHook.HandlerGet() -> Deleting Webhook",
			"Webhook successfully deleted!",
			"",
		)
		debug.ErrorMessage.Print(w)
		return
	}
}

/**
* readData
* Reads data from a Data json object.
**/
func (weatherHook *WeatherHook) ReadData(data interface{}) error {
	m := data.(map[string]interface{})
	weatherHook.Location = m["Location"].(string)
	weatherHook.Timeout = m["Timeout"].(int64)
	weatherHook.URL = m["URL"].(string)
	return nil
}