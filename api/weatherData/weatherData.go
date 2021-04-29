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
	Now struct {
		AirTemperature      float64 `json:"air_temperature"`
		CloudAreaFraction   float64 `json:"cloud_area_fraction"`
		DewPointTemperature	float64 `json:"dew_point_temperature"`
		RelativeHumidity    float64 `json:"relative_humidity"`
		WindFromDirection   float64 `json:"wind_from_direction"`
		WindSpeed           float64 `json:"wind_speed"`
		WindSpeedOfGust     float64 `json:"wind_speed_of_gust"`
		PrecipitationAmount float64 `json:"precipitation_amount"`
	} `json:"now"`
	Today struct {
		Summary                     string  `json:"summary"`
		Confidence                  string  `json:"confidence"`
		AirTemperatureMax           float64 `json:"air_temperature_max"`
		AirTemperatureMin           float64 `json:"air_temperature_min"`
		PrecipitationAmount         float64 `json:"precipitation_amount"`
		PrecipitationAmountMax      float64 `json:"precipitation_amount_max"`
		PrecipitationAmountMin      float64 `json:"precipitation_amount_min"`
		ProbabilityOfPrecipitation	float64 `json:"probability_of_precipitation"`
	} `json:"today"`
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
	weatherData.Now.AirTemperature = yr.Properties.Timeseries[0].Data.Instant.Details.AirTemperature
	weatherData.Now.CloudAreaFraction = yr.Properties.Timeseries[0].Data.Instant.Details.CloudAreaFraction
	weatherData.Now.DewPointTemperature = yr.Properties.Timeseries[0].Data.Instant.Details.DewPointTemperature
	weatherData.Now.RelativeHumidity = yr.Properties.Timeseries[0].Data.Instant.Details.RelativeHumidity
	weatherData.Now.WindFromDirection = yr.Properties.Timeseries[0].Data.Instant.Details.WindFromDirection
	weatherData.Now.WindSpeed = yr.Properties.Timeseries[0].Data.Instant.Details.WindSpeed
	weatherData.Now.WindSpeedOfGust = yr.Properties.Timeseries[0].Data.Instant.Details.WindSpeedOfGust
	weatherData.Now.PrecipitationAmount = yr.Properties.Timeseries[0].Data.Next1Hours.Details.PrecipitationAmount
	
	weatherData.Today.Summary = yr.Properties.Timeseries[0].Data.Next12Hours.Summary.SymbolCode
	weatherData.Today.Confidence = yr.Properties.Timeseries[0].Data.Next12Hours.Summary.SymbolConfidence
	weatherData.Today.AirTemperatureMax = yr.Properties.Timeseries[0].Data.Next6Hours.Details.AirTemperatureMax
	weatherData.Today.AirTemperatureMin = yr.Properties.Timeseries[0].Data.Next6Hours.Details.AirTemperatureMin
	weatherData.Today.PrecipitationAmount = yr.Properties.Timeseries[0].Data.Next6Hours.Details.PrecipitationAmount
	weatherData.Today.PrecipitationAmountMax = yr.Properties.Timeseries[0].Data.Next6Hours.Details.PrecipitationAmountMax
	weatherData.Today.PrecipitationAmountMin = yr.Properties.Timeseries[0].Data.Next6Hours.Details.PrecipitationAmountMin
	weatherData.Today.ProbabilityOfPrecipitation = yr.Properties.Timeseries[0].Data.Next12Hours.Details.ProbabilityOfPrecipitation
	return http.StatusOK, nil
}

func (weatherData *WeatherData) readData(data interface{}) error {
    m := data.(map[string]interface{})
    if data, ok := m["Now"]; ok {
		now := data.(map[string]interface{})
		if field, ok := now["AirTemperature"].(float64); ok {
			weatherData.Now.AirTemperature = field
		}
		if field, ok := now["CloudAreaFraction"].(float64); ok {
			weatherData.Now.CloudAreaFraction = field
		}
		if field, ok := now["DewPointTemperature"].(float64); ok {
			weatherData.Now.DewPointTemperature = field
		}
		if field, ok := now["RelativeHumidity"].(float64); ok {
			weatherData.Now.RelativeHumidity = field
		}
		if field, ok := now["WindFromDirection"].(float64); ok {
			weatherData.Now.WindFromDirection = field
		}
		if field, ok := now["WindSpeed"].(float64); ok {
			weatherData.Now.WindSpeed = field
		}
		if field, ok := now["WindSpeedOfGust"].(float64); ok {
			weatherData.Now.WindSpeedOfGust = field
		}
		if field, ok := now["PrecipitationAmount"].(float64); ok {
			weatherData.Now.PrecipitationAmount = field
		}
    } else {
		return errors.New("getting data from database: Can't find expected fields")
	}
    if data, ok := m["Today"]; ok {
		today := data.(map[string]interface{})
		if field, ok := today["Summary"].(string); ok {
			weatherData.Today.Summary = field
		}
		if field, ok := today["Confidence"].(string); ok {
			weatherData.Today.Confidence = field
		}
		if field, ok := today["AirTemperatureMax"].(float64); ok {
			weatherData.Today.AirTemperatureMax = field
		}
		if field, ok := today["AirTemperatureMin"].(float64); ok {
			weatherData.Today.AirTemperatureMin = field
		}
		if field, ok := today["PrecipitationAmount"].(float64); ok {
			weatherData.Today.PrecipitationAmount = field
		}
		if field, ok := today["PrecipitationAmountMax"].(float64); ok {
			weatherData.Today.PrecipitationAmountMax = field
		}
		if field, ok := today["PrecipitationAmountMin"].(float64); ok {
			weatherData.Today.PrecipitationAmountMin = field
		}
		if field, ok := today["ProbabilityOfPrecipitation"].(float64); ok {
			weatherData.Today.ProbabilityOfPrecipitation = field
		}
    } else {
		return errors.New("getting data from database: Can't find expected fields")
	}
	return nil
}
