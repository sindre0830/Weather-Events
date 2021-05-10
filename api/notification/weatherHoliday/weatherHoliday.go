package weatherHoliday

import (
	"encoding/json"
	"main/api/holidaysData"
	"main/api/notification"
	"main/db"
	"main/debug"
	"net/http"
	"net/url"
	"strings"
)

// WeatherHolidayInput structure, stores information from the user about the webhook
type WeatherHolidayInput struct {
	Holiday 	string `json:"holiday"`
	Location 	string `json:"location"`
	URL 		string `json:"url"`
	Frequency 	string `json:"frequency"`
	Timeout 	int64 `json:"timeout"`
}

// WeatherHoliday structure, stores information about the webhook added to the database
type WeatherHoliday struct {
	ID        	string `json:"id"`
	Date 		string `json:"date"`
	Holiday 	string `json:"holiday"`
	Location 	string `json:"location"`
	URL 		string `json:"url"`
	Frequency 	string `json:"frequency"`
	Timeout 	int64 `json:"timeout"`
}

// Delete a webhook
func (weatherHoliday *WeatherHoliday) Delete(w http.ResponseWriter, r *http.Request) {
	// Parse URL path and ensure that the formatting is correct
	path := strings.Split(r.URL.Path, "/")
	if len(path) != 7 {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"WeatherHoliday.Delete() -> Parsing URL",
			"url validation: either too many or too few arguments in url path",
			"URL format. Remember to add an ID at the end of the path",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	// Get webhook ID
	id := path[len(path)-1]

	err := db.DB.Delete("weatherHoliday", id)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"WeatherHoliday.Delete() -> db.Delete() -> Deleting webhook from the database",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	// Create feedback message to send to client
	var feedback notification.Feedback
	feedback.Update(
		http.StatusNoContent,
		"Webhook successfully deleted",
		id,
	)

	err = feedback.Print(w)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"WeatherHoliday.Delete() -> Feedback.print() -> Sending feedback to client",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
}

// Get a registered webhook
func (weatherHoliday *WeatherHoliday) Get(w http.ResponseWriter, r *http.Request) {
	var holidayMap = make(map[string]interface{})

	// Parse URL path and ensure that the formatting is correct
	path := strings.Split(r.URL.Path, "/")
	if len(path) != 7 {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"WeatherHoliday.View() -> Parsing URL",
			"url validation: either too many or too few arguments in url path",
			"URL format. Remember to add an ID at the end of the path",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	// Get webhook ID
	id := path[len(path)-1]

	if id != ""  {
		// Get webhook from the database if it exists
		data, exists := db.DB.Get("weatherHoliday", id)

		if  !exists {
			debug.ErrorMessage.Update(
				http.StatusBadRequest,
				"WeatherHoliday.Get() -> Database.Get() -> finding document based on ID",
				"getting webhook: can't find id",
				"ID doesn't exist. Expected format: '.../id'. Example: '.../1ab24db3",
			)
			debug.ErrorMessage.Print(w)
			return
		}

		// Parse data into map
		holidayMap = data["Container"].(interface{}).(map[string]interface{})

		// Convert to json
		jsonData, err := json.Marshal(holidayMap)
		if err != nil {
			debug.ErrorMessage.Update(
				http.StatusInternalServerError,
				"WeatherHoliday.Get() -> Converting to json",
				err.Error(),
				"Unknown",
			)
			debug.ErrorMessage.Print(w)
			return
		}

		http.Error(w, string(jsonData), http.StatusOK)
	} else {
		// Get all documents in the collection
		arrData, err := db.DB.GetAll("weatherHoliday")
		if err != nil {
			debug.ErrorMessage.Update(
				http.StatusInternalServerError,
				"WeatherHoliday.Get() -> Database.GetAll() -> Getting all documents",
				err.Error(),
				"Unknown",
			)
			debug.ErrorMessage.Print(w)
			return
		}

		// Add data from the database to the struct
		var weatherHolidays []WeatherHoliday
		for _, rawData := range arrData {
			data := rawData["Container"].(interface{})
			weatherHoliday.readData(data)
			weatherHolidays = append(weatherHolidays, *weatherHoliday)
		}

		// Update header to JSON and set HTTP code
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Send output to user
		err = json.NewEncoder(w).Encode(&weatherHolidays)
		if err != nil {
			debug.ErrorMessage.Update(
				http.StatusInternalServerError,
				"WeatherEvent.GET() -> Sending data to user",
				err.Error(),
				"Unknown",
			)
			debug.ErrorMessage.Print(w)
			return
		}
	}
}

// readData of webhook
func (weatherHoliday *WeatherHoliday) readData(data interface{}) error {
	rawData := data.(map[string]interface{})
	weatherHoliday.ID = rawData["ID"].(string)
	weatherHoliday.Date = rawData["Date"].(string)
	weatherHoliday.Holiday = rawData["Holiday"].(string)
	weatherHoliday.Location = rawData["Location"].(string)
	weatherHoliday.URL = rawData["URL"].(string)
	weatherHoliday.Frequency = rawData["Frequency"].(string)
	weatherHoliday.Timeout = rawData["Timeout"].(int64)
	return nil
}

// POST handles a POST request from the http request.
func (weatherHoliday *WeatherHoliday) Register(w http.ResponseWriter, r *http.Request) {
	var weatherHolidayInput WeatherHolidayInput

	// Decode body into weatherHoliday struct
	err := json.NewDecoder(r.Body).Decode(&weatherHolidayInput)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"WeatherHoliday.Register() -> Decoding body to struct",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	// Check if the URL the user sent is valid
	parsedURL, err := url.ParseRequestURI(weatherHolidayInput.URL)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"WeatherHoliday.Register() -> Checking if URL is valid",
			err.Error(),
			"Not valid URL. Example 'http://google.com/'",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	// Branch if the schema in the URL is incorrect
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"WeatherHoliday.Register() -> Checking if URL is valid",
			"url validation: schema is incorrect, should be 'http' or 'https'",
			"Not valid URL. Example 'http://google.com/'",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	// Check if the timeout sent by the user is valid
	if weatherHolidayInput.Timeout < 15 || weatherHolidayInput.Timeout > 86400 {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"WeatherHoliday.Register() -> Checking if timeout value is valid",
			"timeout validation: value isn't within scope",
			"Timeout value has to be larger then 15 and less then 86400(24 hours) seconds.",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	// Check if the trigger sent by the user is valid
	weatherHolidayInput.Frequency = strings.ToUpper(weatherHolidayInput.Frequency)
	if weatherHolidayInput.Frequency != "EVERY_DAY" && weatherHolidayInput.Frequency != "ON_DATE" {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"WeatherHoliday.Register() -> Checking if trigger value is valid",
			"trigger validation: trigger is not 'EVERY_DAY' or 'ON_DATE'",
			"Not valid trigger. Example 'ON_DATE'",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	// Get a map of all the country's holidays
	var holidaysMap = make(map[string]interface{})
	holidaysMap, status, err := holidaysData.Handler(weatherHolidayInput.Location)
	if err != nil {
		debug.ErrorMessage.Update(
			status,
			"WeatherHoliday.Register() -> holidaysData.Handler() - > Getting information about the country's holidays",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	
	// Make the first letter of each word uppercase to match the format in holidaysMap
	weatherHolidayInput.Holiday = strings.Title(strings.ToLower(weatherHolidayInput.Holiday))

	// Check if the holiday exists in the selected country
	date, ok := holidaysMap[weatherHolidayInput.Holiday]
	if !ok {
		debug.ErrorMessage.Update(
			http.StatusBadRequest,
			"WeatherHoliday.Register() -> Checking if a holiday exists in a country",
			"invalid holiday: the holiday is not valid in the selected country",
			"Not a real holiday. Check your spelling and make sure it is the english name.",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	// Set data to database struct
	weatherHoliday.Date = date.(string)
	weatherHoliday.Holiday = weatherHolidayInput.Holiday
	weatherHoliday.Location = weatherHolidayInput.Location
	weatherHoliday.URL = weatherHolidayInput.URL
	weatherHoliday.Frequency = weatherHolidayInput.Frequency
	weatherHoliday.Timeout = weatherHolidayInput.Timeout

	// Add data to database
	var data db.Data
	data.Container = weatherHoliday

	_, id, err := db.DB.Add("weatherHoliday", "", data)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"WeatherHoliday.Register() -> Database.Add() -> Adding data to database",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	// Create feedback message and print it to the user
	var feedback notification.Feedback
	feedback.Update(
		http.StatusCreated,
		"Webhook successfully created for '" + weatherHoliday.URL + "'",
		id,
	)

	err = feedback.Print(w)
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"WeatherEvent.Register() -> Feedback.print() -> Sending feedback to client",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
}





