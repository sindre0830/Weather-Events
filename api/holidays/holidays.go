package holidays

import (
	"encoding/json"
	"main/api"
	"main/db"
	"main/debug"
	"net/http"
	"strings"
	"time"
)

// Struct for information about one holiday, used when getting data from the API
type Holiday struct {
	Date string `json:"date"`
	Name string `json:"name"`
}

func GetCountryHolidays(w http.ResponseWriter, r *http.Request) {
	var countryHolidays []Holiday
	var newCountryHolidays = make(map[string]interface{})

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

	// Check if country is already stored in the database
	data, exist, err := db.DB.Get("Holidays", location)
	if err != nil && exist {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"WeatherData.Handler() -> Database.get() -> Trying to get data",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	if exist {
		// Convert the data to a map
		m := data.Container.(map[string]interface{})

		// Assign the values to the output map
		for key, elem := range m {
			newCountryHolidays[key] = elem
		}
	} else {
		// Get data from the API and add to the database
		var status int
		countryHolidays, status, err = GetAllHolidays(year, location)
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

		// Put the holidays data in a map where the key is the name and value is the date

		for i := 0; i < len(countryHolidays); i++ {
			newCountryHolidays[countryHolidays[i].Name] = countryHolidays[i].Date
		}

		// Add data to the database
		var data db.Data
		data.Time = time.Now().String()
		data.Container = newCountryHolidays

		_, err = db.DB.Add("Holidays", location, data)
		if err != nil {
			debug.ErrorMessage.Update(
				http.StatusInternalServerError,
				"holidays.GetCountryHolidays() -> Database.Add() -> Adding data to database",
				err.Error(),
				"Unknown",
			)
			debug.ErrorMessage.Print(w)
			return
		}
	}

	// Sending the response to the user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(newCountryHolidays)
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
