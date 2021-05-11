package weatherHook

import (
	"encoding/json"
	"fmt"
	"main/db"
	"main/debug"
	"net/http"
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
	// Start as go routine, else system will hang for the sleep time!
	go weatherHook.Trigger()
	debug.ErrorMessage.Update(
		http.StatusCreated,
		"WeatherHook -> MethodHandler() -> weatherHook.HandlerPost() -> Adding webhook to database.",
		"Webhook successfully added to database! Your ID: "+id,
		"", // We have to add ID here!
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

/**
* readData
* Reads data from a Data struct.
**/
func (weatherHook *WeatherHook) ReadData(data interface{}) error {
	m := data.(map[string]interface{})
	weatherHook.Location = m["Location"].(string)
	weatherHook.Timeout = m["Timeout"].(int64)
	return nil
}

// How to check deletion time?
// Remember to change to hour
/**
* trigger
* This function handles triggering each weatherHook every timeout hours. It's run as a go-routine
**/
func (weatherHook *WeatherHook) Trigger() {
	nextTime := time.Now().Truncate(time.Second) // change to hour!!!!!!!
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

	go weatherHook.Trigger()
}

// Can't put in database due to cyclic import
func StartTrigger(database *db.Database) error {
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
		go temp.Trigger()
	}
	return nil
}
