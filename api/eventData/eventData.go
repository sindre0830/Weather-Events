package eventData

import (
	"errors"
	"main/db"
	"net/http"
	"time"
)

//FirebaseStore -simple struct containing ticketmaster information
type FirebaseStore struct {
	Localdate string `json:"localDate"`
	Name      string `json:"name"`
}

//Handler - Class function will be called and handle all requests and fetches
func (fireBaseStore *FirebaseStore) Handler(eventId string) (int, error) {

	//Check if it is already in firebase
	data, exist := db.DB.Get("Events", eventId)

	if exist { //If in firebase, fetch data from firebase

		//Storing it locally temporarily
		var info FirebaseStore
		err := info.readData(data["Container"].(interface{}))
		if err != nil {
			return http.StatusInternalServerError, err
		}

		//Convert the date string
		layOut := "2006-01-02"
		dateStamp, err := time.Parse(layOut, info.Localdate)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		//Check if the eventdata should be deleted
		overdue := db.CheckIfDateOfEventPassed(dateStamp)

		if overdue {
			//delete the entry in the firebase
			err = db.DB.DeleteEvent(eventId)
			if err != nil {
				return http.StatusInternalServerError, err
			}
		}

		return http.StatusOK, nil

	} else { //If not in firebase, fetch information and store it in firebase

		//Fetch Info:
		var data FirebaseStore
		status, err := data.get(baseURL + eventId + padding + key)
		if err != nil {
			return status, err
		}

		// Add data to the database
		var dataDB db.Data
		dataDB.Container = data

		_, _, err = db.DB.Add("Events", eventId, dataDB)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		return http.StatusOK, nil
	}
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
	data.Name = eventData.Embedded.Venues[0].City.Name
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
	if field, ok := m["Name"].(string); ok {
		data.Name = field
	} else {
		return errors.New("getting data from database: Can't find expected field Name")
	}

	return nil
}
