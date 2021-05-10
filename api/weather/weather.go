package weather

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

// Weather structure stores current and predicted weather data for a day and information about location.
//
// Functionality: Handler, get
type Weather struct {
	Longitude  float64 `json:"longitude"`
	Latitude   float64 `json:"latitude"`
	Location   string  `json:"location"`
	Updated    string  `json:"updated"`
	Timeseries map[string]weatherData.Timeseries `json:"timeseries"`
}

// Handler will handle http request for REST service.
func (weather *Weather) Handler(w http.ResponseWriter, r *http.Request) {
	//parse url and branch if an error occurred
	arrPath := strings.Split(r.URL.Path, "/")
	if len(arrPath) != 6 {
		debug.ErrorMessage.Update(
			http.StatusBadRequest, 
			"Weather.Handler() -> Parsing URL",
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
			"Weather.Handler() -> LocationCoords.Handler() -> Getting location info",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//get all parameters from URL and branch if an error occurred
	arrParam, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError, 
			"Weather.Handler() -> Validating URL parameters",
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
					"Weather.Handler() -> Validating URL parameters",
					err.Error(),
					"Date doesn't match YYYY-MM-DD format. Example: 2021-04-26",
				)
				debug.ErrorMessage.Print(w)
				return
			}
		} else {
			debug.ErrorMessage.Update(
				http.StatusBadRequest, 
				"Weather.Handler() -> Validating URL parameters",
				"url validation: unknown parameter",
				"Unknown",
			)
			debug.ErrorMessage.Print(w)
			return
		}
	}
	//get weather data and branch if an error occurred
	status, err = weather.get(locationCoords.Latitude, locationCoords.Longitude, date)
	if err != nil {
		debug.ErrorMessage.Update(
			status, 
			"Weather.Handler() -> Weather.get() -> Getting weather data",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//set data in structure
	weather.Longitude = locationCoords.Longitude
	weather.Latitude = locationCoords.Latitude
	weather.Location = locationCoords.Address
	//update header to JSON and set HTTP code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//send output to user and branch if an error occured
	err = json.NewEncoder(w).Encode(weather)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError, 
			"Weather.Handler() -> Sending data to user",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
	}
}

// get will get data for structure.
func (weather *Weather) get(lat float64, lon float64, date string) (int, error) {
	//convert coordinates to string
	strLat := fmt.Sprintf("%f", lat)
	strLon := fmt.Sprintf("%f", lon)
	//get weather data and branch if an error occurred
	var weatherDataRange weatherData.WeatherData
	status, err := weatherDataRange.Handler(strLat, strLon)
	if err != nil {
		return status, err
	}
	weather.Timeseries = make(map[string]weatherData.Timeseries)
	//set data in structure
	weather.Updated = weatherDataRange.Updated
	if data, ok := weatherDataRange.Timeseries[date]; ok {
		weather.Timeseries[date] = data
	} else {
		return http.StatusBadRequest, errors.New("invalid date: can't find weather data for inputted date")
	}
	return http.StatusOK, nil
}
