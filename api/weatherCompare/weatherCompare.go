package compare

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/api/geocoords"
	"main/api/weatherData"
	"main/debug"
	"main/fun"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// data structure stores weather data for a location.
type data struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Location  string  `json:"location"`
	Updated   string  `json:"updated"`
	Instant   struct {
		AirTemperature      float64 `json:"air_temperature"`
		CloudAreaFraction   float64 `json:"cloud_area_fraction"`
		DewPointTemperature float64 `json:"dew_point_temperature"`
		RelativeHumidity    float64 `json:"relative_humidity"`
		WindSpeed           float64 `json:"wind_speed"`
		WindSpeedOfGust     float64 `json:"wind_speed_of_gust"`
		PrecipitationAmount float64 `json:"precipitation_amount"`
	} `json:"instant"`
	Predicted struct {
		AirTemperatureMax          float64 `json:"air_temperature_max"`
		AirTemperatureMin          float64 `json:"air_temperature_min"`
		PrecipitationAmount        float64 `json:"precipitation_amount"`
		PrecipitationAmountMax     float64 `json:"precipitation_amount_max"`
		PrecipitationAmountMin     float64 `json:"precipitation_amount_min"`
		ProbabilityOfPrecipitation float64 `json:"probability_of_precipitation"`
	} `json:"predicted"`
}

// locationInfo structure stores all comparison locations information.
type locationInfo struct {
	Location  string
	Longitude float64
	Latitude  float64
}

// WeatherCompare structure stores current and predicted weather data comparisons for different locations.
//
// Functionality: Handler, get
type WeatherCompare struct {
	Longitude  float64 			 `json:"longitude"`
	Latitude   float64 			 `json:"latitude"`
	Location   string  			 `json:"location"`
	Updated    string  			 `json:"updated"`
	Timeseries map[string][]data `json:"timeseries"`
}

// Handler will handle http request for REST service.
func (weatherCompare *WeatherCompare) Handler(w http.ResponseWriter, r *http.Request) {
	//parse url and branch if an error occurred
	arrPath := strings.Split(r.URL.Path, "/")
	if len(arrPath) != 7 {
		debug.ErrorMessage.Update(
			http.StatusBadRequest, 
			"WeatherCompare.Handler() -> Parsing URL",
			"url validation: either too many or too few arguments in url path",
			"URL format. Expected format: '.../place'. Example: '.../oslo'",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	mainLocation := arrPath[5]
	compareLocations := arrPath[6]
	//get locations to compare and branch if there are none
	arrCompareLocations := strings.Split(compareLocations, ";")
	if len(arrCompareLocations) < 1 {
		debug.ErrorMessage.Update(
			http.StatusBadRequest, 
			"WeatherCompare.Handler() -> Getting locations to compare",
			"url validation: not enough locations to compare to",
			"URL format. Expected format: '.../main_place/place1;place2;...'. Example: '.../bergen/oslo;stavanger'",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	var mainLocationCoords geocoords.LocationCoords
	status, err := mainLocationCoords.Handler(mainLocation)
	if err != nil {
		debug.ErrorMessage.Update(
			status, 
			"WeatherCompare.Handler() -> LocationCoords.Handler() -> Getting main location info",
			err.Error(),
			"UnkInstantn",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	weatherCompare.Longitude = mainLocationCoords.Longitude
	weatherCompare.Latitude = mainLocationCoords.Latitude
	weatherCompare.Location = mainLocationCoords.Address
	var arrCoordinates []locationInfo
	for _, location := range arrCompareLocations {
		var locationCoords geocoords.LocationCoords
		status, err := locationCoords.Handler(location)
		if err != nil {
			debug.ErrorMessage.Update(
				status, 
				"WeatherCompare.Handler() -> LocationCoords.Handler() -> Getting comparison location info",
				err.Error(),
				"UnkInstantn",
			)
			debug.ErrorMessage.Print(w)
			return
		}
		var coordinates locationInfo
		coordinates.Longitude = locationCoords.Longitude
		coordinates.Latitude = locationCoords.Latitude
		coordinates.Location = locationCoords.Address
		arrCoordinates = append(arrCoordinates, coordinates)
	}
	//get all parameters from URL and branch if an error occurred
	arrParam, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError, 
			"WeatherCompare.Handler() -> Validating URL parameters",
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
					"WeatherCompare.Handler() -> Validating URL parameters",
					err.Error(),
					"Date doesn't match YYYY-MM-DD format. Example: 2021-04-26",
				)
				debug.ErrorMessage.Print(w)
				return
			}
		} else {
			debug.ErrorMessage.Update(
				http.StatusBadRequest, 
				"WeatherCompare.Handler() -> Validating URL parameters",
				"url validation: unknown parameter",
				"Unknown",
			)
			debug.ErrorMessage.Print(w)
			return
		}
	}
	//get weather data and branch if an error occurred
	status, err = weatherCompare.get(mainLocationCoords.Latitude, mainLocationCoords.Longitude, arrCoordinates, date)
	if err != nil {
		debug.ErrorMessage.Update(
			status, 
			"WeatherCompare.Handler() -> WeatherCompare.get() -> Getting data",
			err.Error(),
			"UnkInstantn",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//update header to JSON and set HTTP code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//send output to user and branch if an error occured
	err = json.NewEncoder(w).Encode(weatherCompare)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError, 
			"WeatherCompare.Handler() -> Sending data to user",
			err.Error(),
			"UnkInstantn",
		)
		debug.ErrorMessage.Print(w)
	}
}

// get will get data for structure.
func (weatherCompare *WeatherCompare) get(lat float64, lon float64, arrCoordinates []locationInfo, date string) (int, error) {
	//convert coordinates to string
	strLat := fmt.Sprintf("%f", lat)
	strLon := fmt.Sprintf("%f", lon)
	//get weather data and branch if an error occurred
	var mainWeatherData weatherData.WeatherData
	status, err := mainWeatherData.Handler(strLat, strLon)
	if err != nil {
		return status, err
	}
	//set data in structure
	weatherCompare.Updated = mainWeatherData.Updated
	//validate date
	if _, ok := mainWeatherData.Timeseries[date]; !ok {
		return http.StatusBadRequest, errors.New("invalid date: can't find weather data for inputted date")
	}
	//get weather data for each comparison location and branch if an error occurred
	weatherCompare.Timeseries = make(map[string][]data)
	var dataRange []data
	for _, coordinates := range arrCoordinates {
		//convert coordinates to string
		strLat := fmt.Sprintf("%f", coordinates.Latitude)
		strLon := fmt.Sprintf("%f", coordinates.Longitude)
		//get weather data and branch if an error occurred
		var weatherData weatherData.WeatherData
		status, err := weatherData.Handler(strLat, strLon)
		if err != nil {
			return status, err
		}
		var data data
		//set data in structure
		data.Longitude = coordinates.Longitude
		data.Latitude = coordinates.Latitude
		data.Location = coordinates.Location
		data.Updated = weatherData.Updated

		data.Instant.AirTemperature = fun.LimitDecimals(weatherData.Timeseries[date].Instant.AirTemperature - mainWeatherData.Timeseries[date].Instant.AirTemperature)
		data.Instant.CloudAreaFraction = fun.LimitDecimals(weatherData.Timeseries[date].Instant.CloudAreaFraction - mainWeatherData.Timeseries[date].Instant.CloudAreaFraction)
		data.Instant.DewPointTemperature = fun.LimitDecimals(weatherData.Timeseries[date].Instant.DewPointTemperature - mainWeatherData.Timeseries[date].Instant.DewPointTemperature)
		data.Instant.RelativeHumidity = fun.LimitDecimals(weatherData.Timeseries[date].Instant.RelativeHumidity - mainWeatherData.Timeseries[date].Instant.RelativeHumidity)
		data.Instant.WindSpeed = fun.LimitDecimals(weatherData.Timeseries[date].Instant.WindSpeed - mainWeatherData.Timeseries[date].Instant.WindSpeed)
		data.Instant.WindSpeedOfGust = fun.LimitDecimals(weatherData.Timeseries[date].Instant.WindSpeedOfGust - mainWeatherData.Timeseries[date].Instant.WindSpeedOfGust)
		data.Instant.PrecipitationAmount = fun.LimitDecimals(weatherData.Timeseries[date].Instant.PrecipitationAmount - mainWeatherData.Timeseries[date].Instant.PrecipitationAmount)

		data.Predicted.AirTemperatureMax = fun.LimitDecimals(weatherData.Timeseries[date].Predicted.AirTemperatureMax - mainWeatherData.Timeseries[date].Predicted.AirTemperatureMax)
		data.Predicted.AirTemperatureMin = fun.LimitDecimals(weatherData.Timeseries[date].Predicted.AirTemperatureMin - mainWeatherData.Timeseries[date].Predicted.AirTemperatureMin)
		data.Predicted.PrecipitationAmount = fun.LimitDecimals(weatherData.Timeseries[date].Predicted.PrecipitationAmount - mainWeatherData.Timeseries[date].Predicted.PrecipitationAmount)
		data.Predicted.PrecipitationAmountMax = fun.LimitDecimals(weatherData.Timeseries[date].Predicted.PrecipitationAmountMax - mainWeatherData.Timeseries[date].Predicted.PrecipitationAmountMax)
		data.Predicted.PrecipitationAmountMin = fun.LimitDecimals(weatherData.Timeseries[date].Predicted.PrecipitationAmountMin - mainWeatherData.Timeseries[date].Predicted.PrecipitationAmountMin)
		data.Predicted.ProbabilityOfPrecipitation = fun.LimitDecimals(weatherData.Timeseries[date].Predicted.ProbabilityOfPrecipitation - mainWeatherData.Timeseries[date].Predicted.ProbabilityOfPrecipitation)
		//append data to array
		dataRange = append(dataRange, data)
	}
	weatherCompare.Timeseries[date] = dataRange
	return http.StatusOK, nil
}
