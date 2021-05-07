package weatherData

import (
	"errors"
	"main/db"
	"net/http"
)

// WeatherData structure stores current and predicted weather data for a day.
//
// Functionality: Handler, get
type WeatherData struct {
	Updated string `json:"updated"`
	Instant struct {
		AirTemperature      float64 `json:"air_temperature"`
		CloudAreaFraction   float64 `json:"cloud_area_fraction"`
		DewPointTemperature	float64 `json:"dew_point_temperature"`
		RelativeHumidity    float64 `json:"relative_humidity"`
		WindFromDirection   float64 `json:"wind_from_direction"`
		WindSpeed           float64 `json:"wind_speed"`
		WindSpeedOfGust     float64 `json:"wind_speed_of_gust"`
		PrecipitationAmount float64 `json:"precipitation_amount"`
	} `json:"instant"`
	Predicted struct {
		Summary                     string  `json:"summary"`
		Confidence                  string  `json:"confidence"`
		AirTemperatureMax           float64 `json:"air_temperature_max"`
		AirTemperatureMin           float64 `json:"air_temperature_min"`
		PrecipitationAmount         float64 `json:"precipitation_amount"`
		PrecipitationAmountMax      float64 `json:"precipitation_amount_max"`
		PrecipitationAmountMin      float64 `json:"precipitation_amount_min"`
		ProbabilityOfPrecipitation	float64 `json:"probability_of_precipitation"`
	} `json:"predicted"`
}

// Handler will handle http request for REST service.
func (weatherData *WeatherData) Handler(lat string, lon string) (int, error) {
	//try to get data from database and branch if an error occurred
	id := lat + "&" + lon
	data, exist, err := db.DB.Get("WeatherData", id)
	if err != nil && exist {
		return http.StatusInternalServerError, err
	}
	//get status on timeframe and branch if an error occurred
	withinTimeframe, err := db.CheckDate(data.Time, 6)
	//check if data is in database and if it's usable then either read data or get new data
	if exist && withinTimeframe {
		if err != nil {
			return http.StatusInternalServerError, err
		}
		err = weatherData.readData(data.Container)
		weatherData.Updated = data.Time
		if err != nil {
			return http.StatusInternalServerError, err
		}
	} else {
		//get data based on coordinates and branch if an error occured
		status, err := weatherData.get(lat, lon)
		if err != nil {
			return status, err
		}
		//send data to database
		var data db.Data
		data.Container = weatherData
		date, err := db.DB.Add("WeatherData", id, data)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		weatherData.Updated = date
	}
	return http.StatusOK, nil
}

// get will get data for structure.
func (weatherData *WeatherData) get(lat string, lon string) (int, error) {
	var yr Yr
	//get weather data from Yr and branch if an error occurred
	status, err := yr.get(lat, lon)
	if err != nil {
		return status, err
	}
	//set data in structure
	weatherData.Instant.AirTemperature = yr.Properties.Timeseries[0].Data.Instant.Details.AirTemperature
	weatherData.Instant.CloudAreaFraction = yr.Properties.Timeseries[0].Data.Instant.Details.CloudAreaFraction
	weatherData.Instant.DewPointTemperature = yr.Properties.Timeseries[0].Data.Instant.Details.DewPointTemperature
	weatherData.Instant.RelativeHumidity = yr.Properties.Timeseries[0].Data.Instant.Details.RelativeHumidity
	weatherData.Instant.WindFromDirection = yr.Properties.Timeseries[0].Data.Instant.Details.WindFromDirection
	weatherData.Instant.WindSpeed = yr.Properties.Timeseries[0].Data.Instant.Details.WindSpeed
	weatherData.Instant.WindSpeedOfGust = yr.Properties.Timeseries[0].Data.Instant.Details.WindSpeedOfGust
	weatherData.Instant.PrecipitationAmount = yr.Properties.Timeseries[0].Data.Next1Hours.Details.PrecipitationAmount
	
	weatherData.Predicted.Summary = yr.Properties.Timeseries[0].Data.Next12Hours.Summary.SymbolCode
	weatherData.Predicted.Confidence = yr.Properties.Timeseries[0].Data.Next12Hours.Summary.SymbolConfidence
	weatherData.Predicted.AirTemperatureMax = yr.Properties.Timeseries[0].Data.Next6Hours.Details.AirTemperatureMax
	weatherData.Predicted.AirTemperatureMin = yr.Properties.Timeseries[0].Data.Next6Hours.Details.AirTemperatureMin
	weatherData.Predicted.PrecipitationAmount = yr.Properties.Timeseries[0].Data.Next6Hours.Details.PrecipitationAmount
	weatherData.Predicted.PrecipitationAmountMax = yr.Properties.Timeseries[0].Data.Next6Hours.Details.PrecipitationAmountMax
	weatherData.Predicted.PrecipitationAmountMin = yr.Properties.Timeseries[0].Data.Next6Hours.Details.PrecipitationAmountMin
	weatherData.Predicted.ProbabilityOfPrecipitation = yr.Properties.Timeseries[0].Data.Next12Hours.Details.ProbabilityOfPrecipitation
	return http.StatusOK, nil
}

func (weatherData *WeatherData) readData(data interface{}) error {
    rawData := data.(map[string]interface{})
    if data, ok := rawData["Instant"]; ok {
		instant := data.(map[string]interface{})
		weatherData.Instant.AirTemperature = instant["AirTemperature"].(float64)
		weatherData.Instant.CloudAreaFraction = instant["CloudAreaFraction"].(float64)
		weatherData.Instant.DewPointTemperature = instant["DewPointTemperature"].(float64)
		weatherData.Instant.RelativeHumidity = instant["RelativeHumidity"].(float64)
		weatherData.Instant.WindFromDirection = instant["WindFromDirection"].(float64)
		weatherData.Instant.WindSpeed = instant["WindSpeed"].(float64)
		weatherData.Instant.WindSpeedOfGust = instant["WindSpeedOfGust"].(float64)
		weatherData.Instant.PrecipitationAmount = instant["PrecipitationAmount"].(float64)
    } else {
		return errors.New("getting data from database: can't find expected fields")
	}
    if data, ok := rawData["Predicted"]; ok {
		predicted := data.(map[string]interface{})
		weatherData.Predicted.Summary = predicted["Summary"].(string)
		weatherData.Predicted.Confidence = predicted["Confidence"].(string)
		weatherData.Predicted.AirTemperatureMax = predicted["AirTemperatureMax"].(float64)
		weatherData.Predicted.AirTemperatureMin = predicted["AirTemperatureMin"].(float64)
		weatherData.Predicted.PrecipitationAmount = predicted["PrecipitationAmount"].(float64)
		weatherData.Predicted.PrecipitationAmountMax = predicted["PrecipitationAmountMax"].(float64)
		weatherData.Predicted.PrecipitationAmountMin = predicted["PrecipitationAmountMin"].(float64)
		weatherData.Predicted.ProbabilityOfPrecipitation = predicted["ProbabilityOfPrecipitation"].(float64)
    } else {
		return errors.New("getting data from database: can't find expected fields")
	}
	return nil
}
