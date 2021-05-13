package holidaysData

import (
	"encoding/json"
	"main/api"
	"main/api/countryData"
	"main/api/geoCoords"
	"main/db"
	"net/http"
	"strings"
	"time"
)

// Handler that gets data about a country's holidaysData from either the API or the database
func Handler(location string) (map[string]interface{}, int, error) {
	var holidaysMap = make(map[string]interface{})

	// Get the coordinates of the location
	var locationCoords geoCoords.LocationCoords
	status, err := locationCoords.Handler(location)
	if err != nil {
		return holidaysMap, http.StatusBadRequest, err
	}

	// Get location's country and format it correctly
	address := strings.Split(locationCoords.Address, ", ")
	country := address[len(address)-1]

	// Get country code
	var countryInfo countryData.Information
	status, err, countryCode := countryInfo.Handler(country)
	if err != nil {
		return holidaysMap, http.StatusBadRequest, err
	}

	// Check if country is already stored in the database
	data, exist := db.DB.Get("Holidays", countryCode)

	if exist {
		// Finds the year the data was saved and the current year
		savedYear := strings.Fields(data["Time"].(string))[2]
		currentYear := strings.Fields(time.Now().Format(time.RFC822))[2]

		// If the years are the same, put data received in the map. If not, get new data from the current year
		if savedYear == currentYear {
			// Convert the data received to a map
			holidaysMap = data["Container"].(interface{}).(map[string]interface{})

			return holidaysMap, http.StatusOK, err
		}
	}

	// Get data from the API and add to the database
	var holidays Holiday
	holidays, status, err = get(countryCode)
	if err != nil {
		return holidaysMap, status, err
	}

	// Put struct data in a map where the key is the name of the holiday and the value is the key
	for i := 0; i < len(holidays); i++ {
		holidaysMap[holidays[i].Name] = holidays[i].Date
	}

	// Add data to the database
	var dataDB db.Data
	dataDB.Container = holidaysMap
	_, _, err = db.DB.Add("Holidays", countryCode, dataDB)
	if err != nil {
		return holidaysMap, http.StatusInternalServerError, err
	}

	return holidaysMap, http.StatusOK, err
}

// get information about all holidaysData in a country
func get(country string) (Holiday, int, error) {
	var holidays Holiday

	// Get the current year as a string
	t := time.Now()
	year := t.String()[:4]

	// Add year and country code as URL path
	url := "https://date.nager.at/api/v2/PublicHolidays/" + year + "/" + country

	// Get data from the request URL
	res, status, err := api.RequestData(url)
	if err != nil {
		return holidays, status, err
	}

	// Unmarshal the response
	err = json.Unmarshal(res, &holidays)
	if err != nil {
		return holidays, http.StatusInternalServerError, err
	}

	return holidays, http.StatusOK, err
}
