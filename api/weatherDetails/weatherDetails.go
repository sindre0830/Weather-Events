package weatherDetails

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/api/geoCoords"
	"main/api/weatherData"
	"main/debug"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Handler will handle http request for REST service.
func (weatherDetails *WeatherDetails) Handler(w http.ResponseWriter, r *http.Request) {
	//parse url and branch if an error occurred
	arrPath := strings.Split(r.URL.Path, "/")
	if len(arrPath) != 6 {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"WeatherDetails.Handler() -> Parsing URL",
			"url validation: either too many or too few arguments in url path",
			"URL format. Expected format: '.../place'. Example: '.../oslo'",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//get location information and branch if an error occurred
	location := arrPath[5]
	var locationCoords geoCoords.LocationCoords
	status, err := locationCoords.Handler(location)
	if err != nil {
		debug.ErrorMessage.Update(
			status,
			"WeatherDetails.Handler() -> LocationCoords.Handler() -> Getting location info",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	weatherDetails.Longitude = locationCoords.Longitude
	weatherDetails.Latitude = locationCoords.Latitude
	weatherDetails.Location = locationCoords.Address
	//get all parameters from URL and branch if an error occurred
	arrParam, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"WeatherDetails.Handler() -> Validating URL parameters",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//branch if any parameters exist
	date := time.Now().Format("2006-01-02")
	if len(arrParam) > 0 {
		//branch if field 'date' exist otherwise return an error
		if targetParameter, ok := arrParam["date"]; ok {
			date = targetParameter[0]
			//validate date
			_, err = time.Parse("2006-01-02", date)
			if err != nil {
				debug.ErrorMessage.Update(
					http.StatusBadRequest,
					"WeatherDetails.Handler() -> Validating URL parameters",
					err.Error(),
					"Date doesn't match YYYY-MM-DD format. Example: 2021-04-26",
				)
				debug.ErrorMessage.Print(w)
				return
			}
		} else {
			debug.ErrorMessage.Update(
				http.StatusBadRequest,
				"WeatherDetails.Handler() -> Validating URL parameters",
				"url validation: unknown parameter",
				"Unknown",
			)
			debug.ErrorMessage.Print(w)
			return
		}
	}
	//get weather data and branch if an error occurred
	status, err = weatherDetails.get(weatherDetails.Latitude, weatherDetails.Longitude, date)
	if err != nil {
		debug.ErrorMessage.Update(
			status,
			"WeatherDetails.Handler() -> WeatherDetails.get() -> Getting weatherDetails data",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//send output to user and branch if an error occured
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(weatherDetails)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"WeatherDetails.Handler() -> Sending data to user",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
	}
}

// get will get data for structure.
func (weatherDetails *WeatherDetails) get(lat float64, lon float64, date string) (int, error) {
	//convert coordinates to string
	strLat := fmt.Sprintf("%f", lat)
	strLon := fmt.Sprintf("%f", lon)
	//get weather data and branch if an error occurred
	var weatherDataRange weatherData.WeatherData
	status, err := weatherDataRange.Handler(strLat, strLon)
	if err != nil {
		return status, err
	}
	//set data in structure and branch if data can't be found
	weatherDetails.Updated = weatherDataRange.Updated
	if data, ok := weatherDataRange.Timeseries[date]; ok {
		weatherDetails.Date = date
		weatherDetails.Data = data
	} else {
		return http.StatusBadRequest, errors.New("invalid date: can't find weather data for inputted date")
	}
	return http.StatusOK, nil
}
