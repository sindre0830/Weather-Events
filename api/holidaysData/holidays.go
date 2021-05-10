package holidaysData

import (
	"encoding/json"
	"main/api"
	"main/db"
	"net/http"
	"strings"
	"time"
)

// Struct for information about one holiday, used when getting data from the API
type Holiday []struct {
	Date string `json:"date"`
	Name string `json:"name"`
}


// Handler that gets data about a country's holidaysData from either the API or the database
func Handler(country string) (map[string]interface{}, int, error) {
	var holidaysMap = make(map[string]interface{})
	var holidays Holiday

	// Check if country is already stored in the database
	data, exist, err := db.DB.Get("Holidays", country)
	if err != nil && exist {
		return holidaysMap, http.StatusInternalServerError, err
	}

	if exist {
		// Finds the year the data was saved and the current year
		savedYear := strings.Fields(data.Time)[2]
		currentYear := strings.Fields(time.Now().Format(time.RFC822))[2]

		// If the years are the same, format the data received from the database. If not, get new data from the current year and add to the database
		if savedYear == currentYear {
			// Convert the data received to a map
			holidaysMap = data.Container.(map[string]interface{})

			// Assign the values to the output map
			for key, elem := range holidaysMap {
				holidaysMap[key] = elem
			}
			return holidaysMap, http.StatusOK, err
		}
	}

	// Get data from the API and add to the database
	status, err := holidays.get(country)
	if err != nil {
		return holidaysMap, status, err
	}

	// Put the holidaysData data in a map where the key is the name and the value is the date
	for i := 0; i < len(holidays); i++ {
		holidaysMap[holidays[i].Name] = holidays[i].Date
	}

	// Add data to the database
	var dataDB db.Data
	dataDB.Container = holidaysMap

	_, _, err = db.DB.Add("Holidays", country, dataDB)
	if err != nil {
		return holidaysMap, http.StatusInternalServerError, err
	}

	return holidaysMap, http.StatusOK, err
}

// get information about all holidaysData in a country
func (holidays *Holiday) get(country string) (int, error) {
	// Get the current year as a string
	t := time.Now()
	year := t.String()[:4]

	// Format the URL
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