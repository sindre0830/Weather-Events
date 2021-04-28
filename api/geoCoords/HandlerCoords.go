package geocoords

import (
	"encoding/json"
	"errors"
	"main/api"
	"main/db"
	"main/debug"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Not sure if we should export this
/*
*	LocationCoords
*	Holds our latitude and longitude data for one location
**/
type LocationCoords struct {
	Latitude	float64	`json:"lat"`
	Longitude	float64 `json:"lon"`
}

var baseURL = "https://us1.locationiq.com/v1/search.php?key="
var key = "pk.d8a67c78822d16869c7a3e8f6d7617af"

/**
*	CoordHandler
*	Accepts a place name, and gets the geological coordinages associated with that place.
*
*	@param	w			-	ResponseWriter we pass our final struct to
*	@param	r			-	Request holding the data passed to us from a user
*
*	@see	getCoords
*	@see	getLocations
*	@see	debug.Debug
**/
func CoordHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	arrPath := strings.Split(r.URL.Path, "/")
	if len(arrPath) != 5 {
		debug.ErrorMessage.Update(
			http.StatusBadRequest, 
			"GeoCoords.Handler() -> Parsing URL",
			"url validation: either too many or too few arguments in url path",
			"URL format. Expected format: '.../place'. Example: '.../oslo'",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	id := strings.ToLower(arrPath[4])
	var locationCoords LocationCoords

	// Check DB fif location data for this location exists
	data, exist, err := db.DB.Get("GeoCoords", id)
	if err != nil && exist {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError, 
			"GeoCoords.HandlerCoords() -> Database.get() -> Trying to get data",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	// We check whether data is deprecated or not.
	// For locations that are not countries/capitals, we don't want to keep our data more than 3 hours.
	withinTimeframe, err := db.CheckDate(data.Time, 3)
	if err != nil && exist {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError, 
			"GeoCoords.HandlerCoords() -> Database.get() -> Trying to get data",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	// If the data is in database and not deprecated, we just read it from firebase and move on
	if exist && withinTimeframe/**withinTimeframe || the location is a country/capital**/ {
		err = readData(&locationCoords, data.Container)
		if err != nil {
			debug.ErrorMessage.Update(
				http.StatusInternalServerError, 
				"GeoCoords.HandlerCoords() -> Database.get() -> reading data",
				err.Error(),
				"Unknown",
			)
			debug.ErrorMessage.Print(w)
			return
		}
	} else {
		// If the location is not stored in firestore, We get the data from the locationiq api
		var locations []map[string]interface{}
		status, err := getLocations(&locations, id)

		if err != nil {
			debug.ErrorMessage.Update(
				status, 
				"GeoCoords.Handler() -> GetCoords.getLocations() -> Getting location data",
				err.Error(),
				"Unknown - ensure that place name is valid and spelled correctly.",
			)
			debug.ErrorMessage.Print(w)
			return
		}
		// We get lat and lon from the first json object in our locations array
		err = getCoords(&locationCoords, locations[0])
		if err != nil {
			debug.ErrorMessage.Update(
				status, 
				"GeoCoords.HandlerCoords() -> WeatherData.get() -> Getting weather data",
				err.Error(),
				"Unknown",
			)
			debug.ErrorMessage.Print(w)
			return
		}
		// Now we send the data to firestore
		var data db.Data
		data.Time = time.Now().String()
		data.Container = locationCoords
		err = db.DB.Add("GeoCoords", id, data)
		if err != nil {
			debug.ErrorMessage.Update(
				http.StatusInternalServerError, 
				"GeoCoords.HandlerCoords() -> Database.Add() -> Adding data to database",
				err.Error(),
				"Unknown",
			)
			debug.ErrorMessage.Print(w)
			return
		}
	}	
	
	// Now that we have our data, we encode and pass it to the user.
	err = json.NewEncoder(w).Encode(locationCoords)
	
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError, 
			"GeoCoords.Handler() -> Sending data to user",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
	}
}

/**
*	getCoords
*	Takes a pointer to a locationCoords struct, and a json object (in map of interfaces form)
*	Finds the relevant fields in the object and puts their data into the struct
*
*	@param	locations		-	Array holding map[string]interface{}
*	@param	coords			-	Struct holding longitude and latitude floats
*	@return	err				-	Interface holding error messages
**/
func getCoords(coords *LocationCoords, location map[string]interface{}) error {
	var err error

	coords.Latitude, err = strconv.ParseFloat(location["lat"].(string), 64) 
	coords.Longitude, err = strconv.ParseFloat(location["lon"].(string), 64) 

	return err
}

/**
*	getLocation
*	Takes a pointer to an array of map[string]interface{} each holding a json object, as well as a location string.
*	Then gets data from the url for that location and puts it into the array.
*
*	@param	locations		-	Array holding map[string]interface{}
*	@param	location		-	String specifying the location for our URL
*	@return	err				-	Interface holding error messages
**/
func getLocations(locations *[]map[string]interface{}, location string) (int, error) {
	url := baseURL + key + "&q=" + location + "&format=json"

	out, status, err := api.RequestData(url)

	if err != nil && status != http.StatusOK {
		return status, err
	}

	// should we use NewDecoder?
	err = json.Unmarshal(out, locations)

	return status, err
}

/**
*	readData
*	Takes a pointer to a LocationCoords struct, and an interface containing the data from Firestore.
*	Reads firestore data into the struct
*
*	@param	coords			-	Pointer to LocationCoords struct we want to fill
*	@param	data			-	interface containing Latitude and Longitude.
*	@return	err				-	Interface holding error messages
**/
func readData(coords *LocationCoords, data interface{}) error {
	m := data.(map[string]interface{})
	if field, ok := m["Latitude"].(float64); ok {
		coords.Latitude = field
	} else {
		return errors.New("getting data from database: Can't find expected fields")
	}
	if field, ok := m["Longitude"].(float64); ok {
		coords.Longitude = field
	} else {
		return errors.New("getting data from database: Can't find expected fields")
	}	
	return nil
}
