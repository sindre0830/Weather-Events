package weatherData

import (
	"errors"
	"main/storage"
	"net/http"
	"time"
)

// Handler will handle http request for REST service.
func (weatherData *WeatherData) Handler(lat string, lon string) (int, error) {
	//get data from database and branch if an error occurred
	id := lat + "&" + lon
	data, exist := storage.Firebase.Get("WeatherData", id)
	//check if data is valid and branch if an error occurred
	var err error
	valid := false
	if date, ok := data["Time"].(string); ok {
		valid, err = storage.CheckDate(date, 3)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}
	//check if data is in database and if it's valid, then either read data or get new data
	if exist && valid {
		//parse data to WeatherData structure and branch if an error occurred
		err := weatherData.readData(data["Container"].(interface{}))
		if err != nil {
			return http.StatusInternalServerError, err
		}
		weatherData.Updated = data["Time"].(string)
	} else {
		//get data based on coordinates and branch if an error occured
		status, err := weatherData.get(lat, lon)
		if err != nil {
			return status, err
		}
		//send data to database and branch if an error occured 
		var data storage.Data
		data.Container = weatherData
		date, _, err := storage.Firebase.Add("WeatherData", id, data)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		weatherData.Updated = date
	}
	return http.StatusOK, nil
}

// get will get data for structure.
func (weatherData *WeatherData) get(lat string, lon string) (int, error) {
	//get weather data from Yr and branch if an error occurred
	var yr Yr
	status, err := yr.get(lat, lon)
	if err != nil {
		return status, err
	}
	//set weather data for all available days
	weatherData.Timeseries = make(map[string]WeatherDataForADay)
	var prevDate string
	for _, elem := range yr.Properties.Timeseries {
		//only add weather data once for each day
		rawDate, _ := time.Parse("2006-01-02T15:04:05Z", elem.Time)
		date := rawDate.Format("2006-01-02")
		if prevDate != date {
			//set data in structure
			var timeEntry WeatherDataForADay
			timeEntry.Instant.AirTemperature = elem.Data.Instant.Details.AirTemperature
			timeEntry.Instant.CloudAreaFraction = elem.Data.Instant.Details.CloudAreaFraction
			timeEntry.Instant.DewPointTemperature = elem.Data.Instant.Details.DewPointTemperature
			timeEntry.Instant.RelativeHumidity = elem.Data.Instant.Details.RelativeHumidity
			timeEntry.Instant.WindFromDirection = elem.Data.Instant.Details.WindFromDirection
			timeEntry.Instant.WindSpeed = elem.Data.Instant.Details.WindSpeed
			timeEntry.Instant.WindSpeedOfGust = elem.Data.Instant.Details.WindSpeedOfGust
			timeEntry.Instant.PrecipitationAmount = elem.Data.Next1Hours.Details.PrecipitationAmount
			timeEntry.Predicted.Summary = elem.Data.Next12Hours.Summary.SymbolCode
			timeEntry.Predicted.Confidence = elem.Data.Next12Hours.Summary.SymbolConfidence
			timeEntry.Predicted.AirTemperatureMax = elem.Data.Next6Hours.Details.AirTemperatureMax
			timeEntry.Predicted.AirTemperatureMin = elem.Data.Next6Hours.Details.AirTemperatureMin
			timeEntry.Predicted.PrecipitationAmount = elem.Data.Next6Hours.Details.PrecipitationAmount
			timeEntry.Predicted.PrecipitationAmountMax = elem.Data.Next6Hours.Details.PrecipitationAmountMax
			timeEntry.Predicted.PrecipitationAmountMin = elem.Data.Next6Hours.Details.PrecipitationAmountMin
			timeEntry.Predicted.ProbabilityOfPrecipitation = elem.Data.Next12Hours.Details.ProbabilityOfPrecipitation
			weatherData.Timeseries[date] = timeEntry
			prevDate = date
		}
	}
	return http.StatusOK, nil
}

// readData parses data from database to WeatherData structure format.
func (weatherData *WeatherData) readData(data interface{}) error {
	//parse nested structure in data and branch if an error occurred
	var timeseries map[string]interface{}
	rawData := data.(map[string]interface{})
	if parsedTimeseries, ok := rawData["Timeseries"].(map[string]interface{}); ok {
		timeseries = parsedTimeseries
	} else {
		return errors.New("parsing data: invalid data structure")
	}
	//parse weather data for all available days
	weatherData.Timeseries = make(map[string]WeatherDataForADay)
	for date, elem := range timeseries {
		data := elem.(map[string]interface{})
		//set data in structure
		var timeEntry WeatherDataForADay
		instant := data["Instant"].(map[string]interface{})
		timeEntry.Instant.AirTemperature = instant["AirTemperature"].(float64)
		timeEntry.Instant.CloudAreaFraction = instant["CloudAreaFraction"].(float64)
		timeEntry.Instant.DewPointTemperature = instant["DewPointTemperature"].(float64)
		timeEntry.Instant.RelativeHumidity = instant["RelativeHumidity"].(float64)
		timeEntry.Instant.WindFromDirection = instant["WindFromDirection"].(float64)
		timeEntry.Instant.WindSpeed = instant["WindSpeed"].(float64)
		timeEntry.Instant.WindSpeedOfGust = instant["WindSpeedOfGust"].(float64)
		timeEntry.Instant.PrecipitationAmount = instant["PrecipitationAmount"].(float64)
		predicted := data["Predicted"].(map[string]interface{})
		timeEntry.Predicted.Summary = predicted["Summary"].(string)
		timeEntry.Predicted.Confidence = predicted["Confidence"].(string)
		timeEntry.Predicted.AirTemperatureMax = predicted["AirTemperatureMax"].(float64)
		timeEntry.Predicted.AirTemperatureMin = predicted["AirTemperatureMin"].(float64)
		timeEntry.Predicted.PrecipitationAmount = predicted["PrecipitationAmount"].(float64)
		timeEntry.Predicted.PrecipitationAmountMax = predicted["PrecipitationAmountMax"].(float64)
		timeEntry.Predicted.PrecipitationAmountMin = predicted["PrecipitationAmountMin"].(float64)
		timeEntry.Predicted.ProbabilityOfPrecipitation = predicted["ProbabilityOfPrecipitation"].(float64)
		weatherData.Timeseries[date] = timeEntry
	}
	return nil
}
