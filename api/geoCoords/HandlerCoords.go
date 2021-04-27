package geocoords

import (
	"encoding/json"
	"main/api"
	"main/debug"
	"net/http"
	"strconv"
	"strings"
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

// Stored list of locations
var Locations = make(map[string]LocationCoords)

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
	if len(arrPath) != 6 {
		debug.ErrorMessage.Update(
			http.StatusBadRequest, 
			"GeoCoords.Handler() -> Parsing URL",
			"url validation: either too many or too few arguments in url path",
			"URL format. Expected format: '.../place'. Example: '.../oslo'",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	placeName := arrPath[4]

	// We get the data from the locationiq api
	var locations []map[string]interface{}
	status, err := getLocations(&locations, placeName)

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
	var locationCoords LocationCoords
	err = getCoords(&locationCoords, locations[0])

	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError, 
			"GeoCoords.Handler() -> GeoCoords.getCoords -> Getting latitude and longitude from json array.",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	locationName :=	strings.ToLower(placeName)

	Locations[locationName] = locationCoords

	err = json.NewEncoder(w).Encode(Locations[locationName])
	
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