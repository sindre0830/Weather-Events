package eventData

import (
	"encoding/json"
	"main/api"
	"net/http"
)

//req -Requests information from the api
func (data *Ticketmaster) req(url string) (int, error) {
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
