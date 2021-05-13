package eventData

import (
	"encoding/json"
	"main/api"
	"net/http"
)

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
