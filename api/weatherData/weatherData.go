package weatherData

import (
	"net/http"
	"strconv"
)

// WeatherData structure stores current and predicted weather data for a day.
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

// get will get data for structure.
func (weatherData *WeatherData) get(lat float64, lon float64) (int, error) {
	var yr Yr
	//convert coordinates to string
	strLat := strconv.FormatFloat(lat, 'f', -1, 64)
	strLon := strconv.FormatFloat(lon, 'f', -1, 64)
	//get weather data from Yr and branch if an error occurred
	status, err := yr.get(strLat, strLon)
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
