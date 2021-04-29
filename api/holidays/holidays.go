package holidays

import (
	"encoding/json"
	"main/api"
	"main/db"
	"net/http"
)

// Struct for information about one holiday, used when getting data from the API
type Holiday []struct {
	Date string `json:"date"`
	Name string `json:"name"`
}


// Gets data about a country's holidays from either the API or the database
func Handler(year string, country string) (map[string]string, int, error) {
	var holidaysMap = make(map[string]string)
	var holidays Holiday

	// Check if country is already stored in the database
	data, exist, err := db.DB.Get("Holidays", country)
	if err != nil && exist {
		return holidaysMap, http.StatusInternalServerError, err
	}

	if exist {
		// Convert the data received to a map
		holidaysMap = data.Container.(map[string]string)

		// Assign the values to the output map
		for key, elem := range holidaysMap {
			holidaysMap[key] = elem
		}
	} else {
		// Get data from the API and add to the database
		status, err := holidays.get(year, country)
		if err != nil {
			return holidaysMap, status, err
		}

		// Put the holidays data in a map where the key is the name and the value is the date
		for i := 0; i < len(holidays); i++ {
			holidaysMap[holidays[i].Name] = holidays[i].Date
		}

		// Add data to the database
		var data db.Data
		data.Container = holidaysMap

		_, err = db.DB.Add("Holidays", country, data)
		if err != nil {
			return holidaysMap, http.StatusInternalServerError, err
		}
	}

	return holidaysMap, http.StatusOK, err
}

// Get information about all holidays in a country
func (holidays *Holiday) get(year string, country string) (int, error) {
	// Slice with holiday structs
	url := "https://date.nager.at/api/v2/PublicHolidays/" + year + "/" + country

	// Gets data from the request URL
	res, status, err := api.RequestData(url)
	if err != nil {
		return status, err
	}

	// Unmarshal the response
	err = json.Unmarshal(res, &holidays)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, err
}