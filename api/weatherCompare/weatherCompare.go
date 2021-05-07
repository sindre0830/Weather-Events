package compare

import (
	"encoding/json"
	"fmt"
	"main/api/geocoords"
	"main/api/weatherData"
	"main/debug"
	"main/fun"
	"net/http"
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
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Location  string  `json:"location"`
	Updated   string  `json:"updated"`
	Data      []data  `json:"data"`
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
	status, err = weatherCompare.get(mainLocationCoords.Latitude, mainLocationCoords.Longitude, arrCoordinates)
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
func (weatherCompare *WeatherCompare) get(lat float64, lon float64, arrCoordinates []locationInfo) (int, error) {
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
	//get weather data for each comparison location and branch if an error occurred
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

		data.Instant.AirTemperature = fun.LimitDecimals(weatherData.Timeseries[time.Now().Format("2006-01-02")].Instant.AirTemperature - mainWeatherData.Timeseries[time.Now().Format("2006-01-02")].Instant.AirTemperature)
		data.Instant.CloudAreaFraction = fun.LimitDecimals(weatherData.Timeseries[time.Now().Format("2006-01-02")].Instant.CloudAreaFraction - mainWeatherData.Timeseries[time.Now().Format("2006-01-02")].Instant.CloudAreaFraction)
		data.Instant.DewPointTemperature = fun.LimitDecimals(weatherData.Timeseries[time.Now().Format("2006-01-02")].Instant.DewPointTemperature - mainWeatherData.Timeseries[time.Now().Format("2006-01-02")].Instant.DewPointTemperature)
		data.Instant.RelativeHumidity = fun.LimitDecimals(weatherData.Timeseries[time.Now().Format("2006-01-02")].Instant.RelativeHumidity - mainWeatherData.Timeseries[time.Now().Format("2006-01-02")].Instant.RelativeHumidity)
		data.Instant.WindSpeed = fun.LimitDecimals(weatherData.Timeseries[time.Now().Format("2006-01-02")].Instant.WindSpeed - mainWeatherData.Timeseries[time.Now().Format("2006-01-02")].Instant.WindSpeed)
		data.Instant.WindSpeedOfGust = fun.LimitDecimals(weatherData.Timeseries[time.Now().Format("2006-01-02")].Instant.WindSpeedOfGust - mainWeatherData.Timeseries[time.Now().Format("2006-01-02")].Instant.WindSpeedOfGust)
		data.Instant.PrecipitationAmount = fun.LimitDecimals(weatherData.Timeseries[time.Now().Format("2006-01-02")].Instant.PrecipitationAmount - mainWeatherData.Timeseries[time.Now().Format("2006-01-02")].Instant.PrecipitationAmount)

		data.Predicted.AirTemperatureMax = fun.LimitDecimals(weatherData.Timeseries[time.Now().Format("2006-01-02")].Predicted.AirTemperatureMax - mainWeatherData.Timeseries[time.Now().Format("2006-01-02")].Predicted.AirTemperatureMax)
		data.Predicted.AirTemperatureMin = fun.LimitDecimals(weatherData.Timeseries[time.Now().Format("2006-01-02")].Predicted.AirTemperatureMin - mainWeatherData.Timeseries[time.Now().Format("2006-01-02")].Predicted.AirTemperatureMin)
		data.Predicted.PrecipitationAmount = fun.LimitDecimals(weatherData.Timeseries[time.Now().Format("2006-01-02")].Predicted.PrecipitationAmount - mainWeatherData.Timeseries[time.Now().Format("2006-01-02")].Predicted.PrecipitationAmount)
		data.Predicted.PrecipitationAmountMax = fun.LimitDecimals(weatherData.Timeseries[time.Now().Format("2006-01-02")].Predicted.PrecipitationAmountMax - mainWeatherData.Timeseries[time.Now().Format("2006-01-02")].Predicted.PrecipitationAmountMax)
		data.Predicted.PrecipitationAmountMin = fun.LimitDecimals(weatherData.Timeseries[time.Now().Format("2006-01-02")].Predicted.PrecipitationAmountMin - mainWeatherData.Timeseries[time.Now().Format("2006-01-02")].Predicted.PrecipitationAmountMin)
		data.Predicted.ProbabilityOfPrecipitation = fun.LimitDecimals(weatherData.Timeseries[time.Now().Format("2006-01-02")].Predicted.ProbabilityOfPrecipitation - mainWeatherData.Timeseries[time.Now().Format("2006-01-02")].Predicted.ProbabilityOfPrecipitation)
		//append data to array
		weatherCompare.Data = append(weatherCompare.Data, data)
	}
	return http.StatusOK, nil
}
