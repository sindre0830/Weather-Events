package weatherData

import (
	"main/db"
	"net/http"
	"time"
)

// Timeseries stores current and predicted weather data for a day
type Timeseries struct {
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

// WeatherData structure stores current and predicted weather data for the next 9 days.
//
// Functionality: Handler, get
type WeatherData struct {
	Updated string `json:"updated"`
	Timeseries map[string]Timeseries `json:"timeseries"`
}

// Handler will handle http request for REST service.
func (weatherData *WeatherData) Handler(lat string, lon string) (int, error) {
	//try to get data from database and branch if an error occurred
	id := lat + "&" + lon
	data, exist := db.DB.Get("WeatherData", id)
	withinTimeframe := false

	// Check if "Time" key exists
	if _, ok := data["Time"].(string); ok {
		//get status on timeframe and branch if an error occurred
		withinTimeframe, _ = db.CheckDate(data["Time"].(string), 3)
	}

	//check if data is in database and if it's usable then either read data or get new data
	if exist && withinTimeframe {
		err := weatherData.readData(data["Container"].(interface{}))
		weatherData.Updated = data["Time"].(string)
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
		date, _, err := db.DB.Add("WeatherData", id, data)
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
	weatherData.Timeseries = make(map[string]Timeseries)
	prevDate := ""
	for _, elem := range yr.Properties.Timeseries {
		var timeEntry Timeseries
		rawDate, _ := time.Parse("2006-01-02T15:04:05Z", elem.Time)
		date := rawDate.Format("2006-01-02")
		if prevDate != date {
			//set data in structure
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

func (weatherData *WeatherData) readData(data interface{}) error {
	weatherData.Timeseries = make(map[string]Timeseries)
    rawData := data.(map[string]interface{})
	timeseries := rawData["Timeseries"].(map[string]interface{})
	for key, elem := range timeseries {
		var timeEntry Timeseries
		data := elem.(map[string]interface{})

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
		weatherData.Timeseries[key] = timeEntry
	}
	return nil
}
