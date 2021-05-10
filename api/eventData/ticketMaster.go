package eventData

import (
	"encoding/json"
	"errors"
	"main/api"
	"main/db"
	"main/debug"
	"net/http"
	"strings"
	"time"
)

//Information -Struct containing all necessary information from ticketmaster
type EventInformation struct {
	Dates struct {
		Start struct {
			Localdate string `json:"localDate"`
		} `json:"start"`
	} `json:"dates"`

	Embedded struct {
		Venues []struct {
			Location struct {
				Longitude string `json:"longitude"`
				Latitude  string `json:"latitude"`
			} `json:"location"`
		} `json:"venues"`
	} `json:"_embedded"`
}

//FirebaseStore -simple struct containing ticketmaster information
type FirebaseStore struct {
	Localdate string `json:"localDate"`
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}

//Constants
var baseURL = "https://app.ticketmaster.com/discovery/v2/events/"
var padding = ".json?apikey="
var key = "ySyIqc6FFKgUIIgzKB5LAOQeGUeU1mot"

//Handler - Class function will be called and handle all requests and fetches
func MethodHandler(w http.ResponseWriter, r *http.Request) {

	//URL parsing
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 7 {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"TicketMaster.Handler() -> Parsing URL",
			"url validation: either too many or too few arguments in url path",
			"URL format. Expected format: '.../event_id'. Example: '.../oslo'",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	id := parts[5]

	//Check if it is already in firebase
	data, exist := db.DB.Get("Events", id)

	if exist { //If in firebase, fetch data from firebase

		//Storing locally
		var info FirebaseStore
		info.readData(data["Container"].(interface{}))

		//Convert the date string
		layOut := "2006-01-02"
		dateStamp, err := time.Parse(layOut, info.Localdate)
		if err != nil {
			debug.ErrorMessage.Update(
				http.StatusInternalServerError,
				"TicketMaster.Handler() -> Converting String Date to time.Time format",
				err.Error(),
				"Unknown",
			)
			debug.ErrorMessage.Print(w)
			return
		}

		//Check if the eventdata should be deleted
		overdue := db.CheckIfDateOfEventPassed(dateStamp)

		if overdue {
			//delete the entry in the firebase
			err = db.DB.Delete(id)
			if err != nil {
				debug.ErrorMessage.Update(
					http.StatusInternalServerError,
					"TicketMaster.Handler() -> Database.Delete() -> Deleting specific event",
					err.Error(),
					"Unknown",
				)
				debug.ErrorMessage.Print(w)
				return
			}
		}

		//update header to JSON and set HTTP code
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		//send output to user and branch if an error occured
		err = json.NewEncoder(w).Encode(info)
		if err != nil {
			debug.ErrorMessage.Update(
				http.StatusInternalServerError,
				"TicketMaster.Handler() -> Sending data to user",
				err.Error(),
				"Unknown",
			)
			debug.ErrorMessage.Print(w)
		}

	} else { //If not in firebase, fetch information and store it in firebase

		//Fetch Info:
		var data FirebaseStore
		status, err := data.get(baseURL + id + padding + key)
		if err != nil {
			debug.ErrorMessage.Update(
				status,
				"TicketMaster.Handler() -> TicketMaster.req() -> Getting data",
				err.Error(),
				"Unknown",
			)
			debug.ErrorMessage.Print(w)
			return
		}

		// Add data to the database
		var dataDB db.Data
		dataDB.Container = data

		_, _, err = db.DB.Add("Events", id, dataDB)
		if err != nil {
			debug.ErrorMessage.Update(
				status,
				"TicketMaster.Handler() -> TicketMaster.req() -> Getting data",
				err.Error(),
				"Unknown",
			)
			debug.ErrorMessage.Print(w)
			return
		}

		//update header to JSON and set HTTP code
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		//send output to user and branch if an error occured
		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			debug.ErrorMessage.Update(
				http.StatusInternalServerError,
				"WeatherCompare.Handler() -> Sending data to user",
				err.Error(),
				"Unknown",
			)
			debug.ErrorMessage.Print(w)
		}
	}
}

//req -Requests information from the api
func (data *EventInformation) req(url string) (int, error) {
	output, status, jsonErr := api.RequestData(url)

	if jsonErr != nil {
		return status, jsonErr
	}

	jsonErr = json.Unmarshal(output, &data)
	if jsonErr != nil {
		return http.StatusInternalServerError, jsonErr
	}

	return http.StatusOK, nil
}

//get -Requests information from the api, going through the more complex struct
func (data *FirebaseStore) get(url string) (int, error) {
	var eventData EventInformation
	status, err := eventData.req(url) //Get information into the other struct
	if err != nil {
		return http.StatusInternalServerError, err
	}

	//For ease store that information in a more readable format
	data.Localdate = eventData.Dates.Start.Localdate
	data.Latitude = eventData.Embedded.Venues[0].Location.Latitude
	data.Longitude = eventData.Embedded.Venues[0].Location.Longitude
	return status, nil
}

//readData -Stores information from interface to the struct
func (data *FirebaseStore) readData(storage interface{}) error {
	m := storage.(map[string]interface{})
	if field, ok := m["Localdate"].(string); ok {
		data.Localdate = field
	} else {
		return errors.New("getting data from database: Can't find expected field Localdate")
	}
	if field, ok := m["Latitude"].(string); ok {
		data.Latitude = field
	} else {
		return errors.New("getting data from database: Can't find expected field Latitude")
	}
	if field, ok := m["Longitude"].(string); ok {
		data.Longitude = field
	} else {
		return errors.New("getting data from database: Can't find expected field Latitude")
	}
	return nil
}
