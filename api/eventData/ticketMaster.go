package eventData

import (
	"encoding/json"
	"main/api"
	"net/http"
)

//Information -Struct containing all necessary information from ticketmaster
type EventInformation2 struct {
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

type EventInformation struct {
	Dates struct {
		Start struct {
			Localdate string `json:"localDate"`
		} `json:"start"`
	} `json:"dates"`

	Embedded struct {
		Venues []struct {
			City struct {
				Name string `json:"name"`
			} `json:"city"`
		} `json:"venues"`
	} `json:"_embedded"`
}

//Constants
var baseURL = "https://app.ticketmaster.com/discovery/v2/events/"
var padding = ".json?apikey="
var key = "ySyIqc6FFKgUIIgzKB5LAOQeGUeU1mot"

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
