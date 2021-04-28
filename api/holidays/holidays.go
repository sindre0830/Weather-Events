package holidays

import (
	"encoding/json"
	"main/api"
	"main/debug"
	"net/http"
	"strings"
)

// Struct for information about one holiday, used when getting data from the API
type Holiday struct {
	Date        string   `json:"date"`
	LocalName   string   `json:"localName"`
	Name        string   `json:"name"`
	CountryCode string   `json:"countryCode"`
	Fixed       bool     `json:"fixed"`
	Global      bool     `json:"global"`
	Counties    []string `json:"counties"`
	LaunchYear  int      `json:"launchYear"`
	Type        string   `json:"string"`
}

func GetCountryHolidays(w http.ResponseWriter, r *http.Request) {
	var countryHolidays []Holiday

	// Parsing variables from URL path
	path := strings.Split(r.URL.Path, "/")
	// Check if the path is correctly formed
	if len(path) != 6 {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"holidays.GetCountryHolidays() -> Parsing URL",
			"url validation: either too many or too few arguments in url path",
			"URL format. Expected format: '.../year/location'. Example: '.../2021/NO'",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	year := path[4]
	location := path[5]

	// Get the country's holiday
	countryHolidays, status, err := GetAllHolidays(year, location)
	if err != nil {
		debug.ErrorMessage.Update(
			status,
			"holidays.GetCountryHolidays() -> holidays.GetAllHolidays() -> Getting holidays",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	// Sending the response to the user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(countryHolidays)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"holidays.GetCountryHolidays() -> Sending data to user",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)

	}
}

// Get information about all holidays in a country
func GetAllHolidays(year string, countryCode string) ([]Holiday, int, error) {
	// Slice with holiday structs
	var countryHolidays []Holiday
	url := "https://date.nager.at/api/v2/PublicHolidays/" + year + "/" + countryCode

	// Gets data from the request URL
	res, status, err := api.RequestData(url)
	if err != nil {
		return countryHolidays, status, err
	}

	// Unmarshal the response
	err = json.Unmarshal(res, &countryHolidays)
	if err != nil {
		return countryHolidays, http.StatusInternalServerError, err
	}

	return countryHolidays, http.StatusOK, err
}
