package weatherData

import (
	"encoding/json"
	"main/api"
	"net/http"
)

// Data structure stores current and all predicted weather data for a day.
type Data struct {
	Instant struct {
		Details struct {
			AirPressureAtSeaLevel    float64 `json:"air_pressure_at_sea_level"`
			AirTemperature           float64 `json:"air_temperature"`
			CloudAreaFraction        float64 `json:"cloud_area_fraction"`
			CloudAreaFractionHigh    float64 `json:"cloud_area_fraction_high"`
			CloudAreaFractionLow     float64 `json:"cloud_area_fraction_low"`
			CloudAreaFractionMedium  float64 `json:"cloud_area_fraction_medium"`
			DewPointTemperature      float64 `json:"dew_point_temperature"`
			FogAreaFraction          float64 `json:"fog_area_fraction"`
			RelativeHumidity         float64 `json:"relative_humidity"`
			UltravioletIndexClearSky float64 `json:"ultraviolet_index_clear_sky"`
			WindFromDirection        float64 `json:"wind_from_direction"`
			WindSpeed                float64 `json:"wind_speed"`
			WindSpeedOfGust          float64 `json:"wind_speed_of_gust"`
		} `json:"details"`
	} `json:"instant"`
	Next12Hours struct {
		Summary struct {
			SymbolCode       string `json:"symbol_code"`
			SymbolConfidence string `json:"symbol_confidence"`
		} `json:"summary"`
		Details struct {
			ProbabilityOfPrecipitation float64 `json:"probability_of_precipitation"`
		} `json:"details"`
	} `json:"next_12_hours"`
	Next1Hours struct {
		Summary struct {
			SymbolCode string `json:"symbol_code"`
		} `json:"summary"`
		Details struct {
			PrecipitationAmount        float64 `json:"precipitation_amount"`
			PrecipitationAmountMax     float64 `json:"precipitation_amount_max"`
			PrecipitationAmountMin     float64 `json:"precipitation_amount_min"`
			ProbabilityOfPrecipitation float64 `json:"probability_of_precipitation"`
			ProbabilityOfThunder       float64 `json:"probability_of_thunder"`
		} `json:"details"`
	} `json:"next_1_hours"`
	Next6Hours struct {
		Summary struct {
			SymbolCode string `json:"symbol_code"`
		} `json:"summary"`
		Details struct {
			AirTemperatureMax          float64 `json:"air_temperature_max"`
			AirTemperatureMin          float64 `json:"air_temperature_min"`
			PrecipitationAmount        float64 `json:"precipitation_amount"`
			PrecipitationAmountMax     float64 `json:"precipitation_amount_max"`
			PrecipitationAmountMin     float64 `json:"precipitation_amount_min"`
			ProbabilityOfPrecipitation float64 `json:"probability_of_precipitation"`
		} `json:"details"`
	} `json:"next_6_hours"`
}

// Yr structure stores weather data for the next 10 days.
//
// Functionality: get, req
type Yr struct {
	Properties struct {
		Timeseries []struct {
			Data Data `json:"data"`
		} `json:"timeseries"`
	} `json:"properties"`
}

// get will get data for structure.
func (yr *Yr) get(lat string, lon string) (int, error) {
	url := "https://api.met.no/weatherapi/locationforecast/2.0/complete?lat=" + lat + "&lon=" + lon
	//gets json output from API and branch if an error occurred
	status, err := yr.req(url)
	if err != nil {
		return status, err
	}
	return http.StatusOK, nil
}

// req will request data from API.
func (yr *Yr) req(url string) (int, error) {
	//gets raw data from API and branch if an error occurred
	output, status, err := api.RequestData(url)
	if err != nil {
		return status, err
	}
	//convert raw data to JSON and branch if an error occurred
	err = json.Unmarshal(output, &yr)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
