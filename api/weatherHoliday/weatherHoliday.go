package weatherHoliday

import (
	"encoding/json"
	"fmt"
	"main/api/countryData"
	"main/api/geocoords"
	"main/api/holidaysData"
	"main/debug"
	"net/http"
	"strings"
)

type WeatherHoliday struct {
	Holiday string `json:"holiday"`
	Location string `json:"location"`
	Frequency int `json:"frequency"`
}

// Handler for the weather holiday webhook endpoint
func (weatherHoliday *WeatherHoliday) Handler(w http.ResponseWriter, r *http.Request) {
	// Decode body into struct
	err := json.NewDecoder(r.Body).Decode(&weatherHoliday)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"WeatherHoliday.Handler() -> Decoding body to struct",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	// Get the geocoords of the location
	var locationCoords geocoords.LocationCoords
	status, err := locationCoords.Handler(weatherHoliday.Location)
	if err != nil {
		debug.ErrorMessage.Update(
			status,
			"WeatherHoliday.Handler() -> LocationCoords.Handler() -> Getting location info",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	// Get country and format it correctly
	address := strings.Split(locationCoords.Address, ",")	// TODO: ", "
	country := address[len(address)-1]
	country = country[1:]

	// Get country code
	countryCode, status, err := getCode(country)
	if err != nil {
		debug.ErrorMessage.Update(
			status,
			"WeatherHoliday.Handler() -> WeatherHoliday.getCode() -> Getting country code",
			err.Error(),
			"Selected country is not valid",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	fmt.Println(countryCode)

	// Get the country's holidays
	var holidaysMap = make(map[string]interface{})
	holidaysMap, status, err = holidaysData.Handler(countryCode)
	if err != nil {
		debug.ErrorMessage.Update(
			status,
			"WeatherHoliday.Handler() -> holidaysData.Handler() - > Getting information about the country's holidays",
			err.Error(),
			"Unknown",
		)
	}

	// Check if the holiday exists in the selected country
	_, ok := holidaysMap[weatherHoliday.Holiday]

	if !ok {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"WeatherHoliday.Handler() -> Checking if holidays exists in country",
			err.Error(),
			"The selected holiday is not valid",
		)
		debug.ErrorMessage.Print(w)
		return
	}

}

// getCode - Get country's alpha code
func getCode(countryName string) (string, int, error) {
	var countryInfo countryData.Information
	status, err, countryCode := countryInfo.Handler(countryName)
	if err != nil {
		return "", status, err
	}

	return countryCode, status, err
}

