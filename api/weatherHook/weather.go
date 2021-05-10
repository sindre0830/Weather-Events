package weatherHook

import (
	"encoding/json"
	"fmt"
	"main/db"
	"main/debug"
	"net/http"
	"net/url"
	"time"
)

type WeatherHook struct {
	Location 	string    `json:"location"`
	Timeout		int64	  `json:"timeout"`
}

/**
* HandlerPost
* Handles POST method requests from client.
**/
func (weatherHook *WeatherHook) HandlerPost(w http.ResponseWriter, r *http.Request) {
	// turn json object into struct
	err := json.NewDecoder(r.Body).Decode(&weatherHook)
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
	// add it to firebase
	// How do we want to handle ID, getting/passing to user? Currently getting but not checking dupes.
	var data db.Data
	data.Container = weatherHook
	_, id, err := db.DB.Add("notification", "", data)
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

	weatherHook.trigger()
	debug.ErrorMessage.Update(
		http.StatusCreated, 
		"WeatherHook -> MethodHandler() -> weatherHook.HandlerPost() -> Adding webhook to database.",
		"Webhook successfully added to database! Your ID: " + id,
		"",				// We have to add ID here!
	)
	debug.ErrorMessage.Print(w)
	return
}


/**
* HandlerGet
* Handles GET method requests from clients
* This method should post the webhook itself to the client.
**/
func (weatherHook *WeatherHook) HandlerGet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get func")
	w.Header().Set("Content-Type", "application/json")
	// get query with ID
	params, _ := url.ParseQuery(r.URL.RawQuery)
	id := params["id"][0]

	// check for ID in firestore - DB function
	data, exist, err := db.DB.Get("notification", id)		// all hooks in one db?
	if err != nil && exist {
		debug.ErrorMessage.Update(
			http.StatusNotFound, 
			"WeatherHook -> MethodHandler() -> weatherHook.HandlerGet() -> Getting data from firestore.",
			"Database error: Failed when getting data from Firestore!",
			"",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	// Extract from data
	if exist {
		err = weatherHook.readData(data.Container)
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
		err = json.NewEncoder(w).Encode(weatherHook)
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
	fmt.Println("Delete func")
	w.Header().Set("Content-Type", "application/json")
	// get query with ID
	params, _ := url.ParseQuery(r.URL.RawQuery)
	id := params["id"][0]
	// delete
	err := db.DB.Delete(id)

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
* Reads data from a Data struct.
**/
func (weatherHook *WeatherHook) readData(data interface{}) error {
	m := data.(map[string]interface{})
	weatherHook.Location = m["Location"].(string)
	weatherHook.Timeout = int64(m["Timeout"].(float64))
	return nil
}

// Currently only runs til the program goes down
func (weatherHook *WeatherHook) trigger () {
	nextTime := time.Now().Truncate(time.Second)	// change to hour
	nextTime = nextTime.Add(time.Duration(weatherHook.Timeout) * time.Second)
	time.Sleep(time.Until(nextTime))

	var url = "/weather-rest/v1/weather/location/" + weatherHook.Location
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
	fmt.Printf("\nError when creating new POST request.\nRaw error: %v\n", err.Error())
	} else {
		fmt.Printf("\nPassed webhook to user!")
	}
	
	go weatherHook.trigger()
}