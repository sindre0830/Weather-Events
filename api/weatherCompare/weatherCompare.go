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
)

type Data struct {
	Location  string  `json:"location"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latiude"`
	Now       struct {
		AirTemperature      float64 `json:"air_temperature"`
		CloudAreaFraction   float64 `json:"cloud_area_fraction"`
		DewPointTemperature float64 `json:"dew_point_temperature"`
		RelativeHumidity    float64 `json:"relative_humidity"`
		WindSpeed           float64 `json:"wind_speed"`
		WindSpeedOfGust     float64 `json:"wind_speed_of_gust"`
		PrecipitationAmount float64 `json:"precipitation_amount"`
	} `json:"now"`
	Today struct {
		AirTemperatureMax          float64 `json:"air_temperature_max"`
		AirTemperatureMin          float64 `json:"air_temperature_min"`
		PrecipitationAmount        float64 `json:"precipitation_amount"`
		PrecipitationAmountMax     float64 `json:"precipitation_amount_max"`
		PrecipitationAmountMin     float64 `json:"precipitation_amount_min"`
		ProbabilityOfPrecipitation float64 `json:"probability_of_precipitation"`
	} `json:"today"`
}

// WeatherCompare structure stores current and predicted weather data comparisons for different locations.
type WeatherCompare struct {
	Updated      string `json:"updated"`
	MainLocation string `json:"main_location"`
	Data         []Data `json:"data"`
}

type LocationInfo struct {
	Location  string  `json:"location"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latiude"`
}

func (weatherCompare *WeatherCompare) Handler(w http.ResponseWriter, r *http.Request) {
	//parse url and branch if an error occurred
	arrPath := strings.Split(r.URL.Path, "/")
	if len(arrPath) != 7 {
		debug.ErrorMessage.Update(
			http.StatusBadRequest, 
			"WeatherCompare.Handler() -> Parsing URL",
			"url validation: either too many or too few arguments in url path",
			"URL format. Expected format: '.../place?fields=filter1;filter2...'. Example: '.../oslo?fields=Data'",
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
			"URL format. Expected format: '...main_place/place1;place2;...'. Example: '...bergen/oslo;stavanger'",
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
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	weatherCompare.MainLocation = mainLocationCoords.Address
	var arrCoordinates []LocationInfo
	for _, location := range arrCompareLocations {
		var locationCoords geocoords.LocationCoords
		status, err := locationCoords.Handler(location)
		if err != nil {
			debug.ErrorMessage.Update(
				status, 
				"WeatherCompare.Handler() -> LocationCoords.Handler() -> Getting comparison location info",
				err.Error(),
				"Unknown",
			)
			debug.ErrorMessage.Print(w)
			return
		}
		var coordinates LocationInfo
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
			"Unknown",
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
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
	}
}

func (weatherCompare *WeatherCompare) get(lat float64, lon float64, arrCoordinates []LocationInfo) (int, error) {
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
		var data Data
		//set data in structure
		data.Location = coordinates.Location
		data.Longitude = coordinates.Longitude
		data.Latitude = coordinates.Latitude

		data.Now.AirTemperature = fun.LimitDecimals(weatherData.Now.AirTemperature - mainWeatherData.Now.AirTemperature)
		data.Now.CloudAreaFraction = fun.LimitDecimals(weatherData.Now.CloudAreaFraction - mainWeatherData.Now.CloudAreaFraction)
		data.Now.DewPointTemperature = fun.LimitDecimals(weatherData.Now.DewPointTemperature - mainWeatherData.Now.DewPointTemperature)
		data.Now.RelativeHumidity = fun.LimitDecimals(weatherData.Now.RelativeHumidity - mainWeatherData.Now.RelativeHumidity)
		data.Now.WindSpeed = fun.LimitDecimals(weatherData.Now.WindSpeed - mainWeatherData.Now.WindSpeed)
		data.Now.WindSpeedOfGust = fun.LimitDecimals(weatherData.Now.WindSpeedOfGust - mainWeatherData.Now.WindSpeedOfGust)
		data.Now.PrecipitationAmount = fun.LimitDecimals(weatherData.Now.PrecipitationAmount - mainWeatherData.Now.PrecipitationAmount)

		data.Today.AirTemperatureMax = fun.LimitDecimals(weatherData.Today.AirTemperatureMax - mainWeatherData.Today.AirTemperatureMax)
		data.Today.AirTemperatureMin = fun.LimitDecimals(weatherData.Today.AirTemperatureMin - mainWeatherData.Today.AirTemperatureMin)
		data.Today.PrecipitationAmount = fun.LimitDecimals(weatherData.Today.PrecipitationAmount - mainWeatherData.Today.PrecipitationAmount)
		data.Today.PrecipitationAmountMax = fun.LimitDecimals(weatherData.Today.PrecipitationAmountMax - mainWeatherData.Today.PrecipitationAmountMax)
		data.Today.PrecipitationAmountMin = fun.LimitDecimals(weatherData.Today.PrecipitationAmountMin - mainWeatherData.Today.PrecipitationAmountMin)
		data.Today.ProbabilityOfPrecipitation = fun.LimitDecimals(weatherData.Today.ProbabilityOfPrecipitation - mainWeatherData.Today.ProbabilityOfPrecipitation)
		//append data to array
		weatherCompare.Data = append(weatherCompare.Data, data)
	}
	return http.StatusOK, nil
}
