package weatherData

import (
	"encoding/json"
	"main/api"
	"main/dict"
	"net/http"
)

// get will get data for structure.
func (yr *Yr) get(lat string, lon string) (int, error) {
	url := dict.GetYrURL(lat, lon)
	//gets json output from API and branch if an error occurred
	status, err := yr.req(url)
	if err != nil {
		return status, err
	}
	return http.StatusOK, nil
}

// req will request data from API.
func (yr *Yr) req(url string) (int, error) {
	//gets raw data from API and branch if an error occurred
	output, status, err := api.RequestData(url)
	if err != nil {
		return status, err
	}
	//convert raw data to JSON and branch if an error occurred
	err = json.Unmarshal(output, &yr)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
