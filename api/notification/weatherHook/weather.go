package weatherHook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"main/api/notification"
	"main/api/weather"
	"main/db"
	"main/debug"
	"main/dict"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"google.golang.org/api/iterator"
)

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
	//validate parameters and branch if an error occurred
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
	//check if timeout is valid and return an error if it isn't - timeout in hours
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
	// How do we want to handle ID, getting/passing to user? Currently getting but not checking dupes.
	var data db.Data
	data.Container = weatherHook
	_, id, err := db.DB.Add(hookdb, "", data)
	// Return ID if successful, otherwise error
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
	//start loop
	go weatherHook.callLoop()
	//create feedback message to send to client and branch if an error occurred
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
	// // Start as go routine, else system will hang for the sleep time!
	// go weatherHook.Trigger()
	// debug.ErrorMessage.Update(
	// 	http.StatusCreated,
	// 	"WeatherHook -> MethodHandler() -> weatherHook.HandlerPost() -> Adding webhook to database.",
	// 	"Webhook successfully added to database! Your ID: " + id,
	// 	"",				// We have to add ID here!
	// )
	// debug.ErrorMessage.Print(w)
	// return
}

/**
* HandlerGet
* Handles GET method requests from clients
* This method should post the webhook itself to the client.
**/
func (weatherHook *WeatherHook) HandlerGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// get query with ID
	params, _ := url.ParseQuery(r.URL.RawQuery)
	id := params["id"][0]

	// check for ID in firestore - DB function
	data, exist := db.DB.Get(hookdb, id) // all hooks in one db?
	// Extract from data
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
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(&weatherHook)
		//send output to user and branch if an error occured
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
	params, _ := url.ParseQuery(r.URL.RawQuery)
	id := params["id"][0]
	// delete
	err := db.DB.Delete(hookdb, id)

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

func (weatherHook *WeatherHook) callLoop() {
	_, exist := db.DB.Get("weatherHook", weatherHook.ID)
	if !exist {
		return
	}
	nextTime := time.Now().Truncate(time.Second)
	nextTime = nextTime.Add(time.Duration(weatherHook.Timeout) * time.Second) // Change to hour!
	time.Sleep(time.Until(nextTime))
	url := dict.GetWeatherURL(weatherHook.Location, "")
	var weather weather.Weather
	//create new GET request and branch if an error occurred
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf(
			"%v {\n\tError when creating HTTP request to Weather.Handler().\n\tRaw error: %v\n}\n",
			time.Now().Format("2006-01-02 15:04:05"), err.Error(),
		)
		go weatherHook.callLoop()
	}
	//call the policy handler and branch if the status code is not OK
	//this stops timed out request being sent to the webhook
	recorder := httptest.NewRecorder()
	weather.Handler(recorder, req)
	if recorder.Result().StatusCode != http.StatusOK {
		fmt.Printf(
			"%v {\n\tError when creating HTTP request to Weather.Handler().\n\tStatus code: %v\n}\n",
			time.Now().Format("2006-01-02 15:04:05"), recorder.Result().StatusCode,
		)
		go weatherHook.callLoop()
	}
	//convert from structure to bytes and branch if an error occurred
	output, err := json.Marshal(weather)
	if err != nil {
		fmt.Printf(
			"%v {\n\tError when parsing Weather structure.\n\tRaw error: %v\n}\n",
			time.Now().Format("2006-01-02 15:04:05"), err.Error(),
		)
		go weatherHook.callLoop()
	}
	//create new POST request and branch if an error occurred
	req, err = http.NewRequest(http.MethodPost, weatherHook.URL, bytes.NewBuffer(output))
	if err != nil {
		fmt.Printf(
			"%v {\n\tError when creating new POST request.\n\tRaw error: %v\n}\n",
			time.Now().Format("2006-01-02 15:04:05"), err.Error(),
		)
		go weatherHook.callLoop()
	}
	//hash structure and branch if an error occurred
	mac := hmac.New(sha256.New, dict.Secret)
	_, err = mac.Write([]byte(output))
	if err != nil {
		fmt.Printf(
			"%v {\n\tError when hashing content before POST request.\n\tRaw error: %v\n}\n",
			time.Now().Format("2006-01-02 15:04:05"), err.Error(),
		)
		go weatherHook.callLoop()
	}
	//convert hashed structure to string and add to header
	req.Header.Add("Signature", hex.EncodeToString(mac.Sum(nil)))
	//update header to JSON
	req.Header.Set("Content-Type", "application/json")
	//send request to client and branch if an error occured
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf(
			"%v {\n\tError when sending HTTP content to webhook.\n\tRaw error: %v\n}\n",
			time.Now().Format("2006-01-02 15:04:05"), err.Error(),
		)
		go weatherHook.callLoop()
	}
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusServiceUnavailable {
		fmt.Printf(
			"%v {\n\tWebhook URL is not valid. Deleting webhook...\n\tStatus code: %v\n}\n",
			time.Now().Format("2006-01-02 15:04:05"), res.StatusCode,
		)
		db.DB.Delete("weatherEvent", weatherHook.ID)
		return
	}
	go weatherHook.callLoop()
}

/**
* readData
* Reads data from a Data struct.
**/
func (weatherHook *WeatherHook) ReadData(data interface{}) error {
	m := data.(map[string]interface{})
	weatherHook.Location = m["Location"].(string)
	weatherHook.Timeout = m["Timeout"].(int64)
	weatherHook.URL = m["URL"].(string)
	return nil
}

// Can't put in database due to cyclic import
func StartCall(database *db.Database) error {
	iter := database.Client.Collection("weatherHookDB").Documents(database.Ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		// Create dummy variables
		var temp WeatherHook
		var dbMap = doc.Data()
		// Extract data from iterator
		err = temp.ReadData(dbMap["Container"].(interface{}))
		if err != nil {
			return err
		}
		// Start as go routine, else system will hang for the sleep time!
		go temp.callLoop()
	}
	return nil
}
