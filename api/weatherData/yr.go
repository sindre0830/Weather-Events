package weatherData

import (
	"encoding/json"
	"main/api"
	"net/http"
)

// Yr structure stores weather data within scope for a location.
//
// Functionality: Get, req
type Yr struct {
	Properties struct {
		Meta       interface{}   `json:"meta"`
		Timeseries []interface{} `json:"timeseries"`
	} `json:"properties"`
}

// Get will get data for structure.
func (yr *Yr) Get(lat string, lon string) (int, error) {
	url := "https://api.met.no/weatherapi/locationforecast/2.0/complete?lat=" + lat + "&lon=" + lon
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
