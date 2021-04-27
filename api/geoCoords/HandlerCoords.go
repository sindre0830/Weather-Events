package geocoords

import (
	"encoding/json"
	"main/api"
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
	parts := strings.Split(r.URL.Path, "/")
	// Implement this check once error setup is complete
	// if len(parts) != 5 {
	// 	outErr := api.MakeError("In function CoordHandler, expecting format weather-rest/v1/geocoord/PLACE", http.StatusBadRequest)

	// 	json.NewEncoder(w).Encode(outErr)
	// 	return
	// }
	// place := parseName(parts[4])  // Not sure we need parseName function? Will check later.
	placeName := parts[4]

	// We get the data from the locationiq api
	var locations []map[string]interface{}
	err := getLocations(&locations, placeName)

	if err != nil {
		err = json.NewEncoder(w).Encode(err)
	}

	// We pass the first json object in our locations array
	var locationCoords LocationCoords
	err = getCoords(&locationCoords, locations[0])

	locationName :=	strings.ToLower(placeName)

	Locations[locationName] = locationCoords

	// Again, fix up once error handling is clarified
	if err != nil {
		err = json.NewEncoder(w).Encode(err)
	}

	err = json.NewEncoder(w).Encode(Locations[locationName])
	
	// Again, fix up once error handling is clarified
	if err != nil {
		err = json.NewEncoder(w).Encode(err)
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
func getLocations(locations *[]map[string]interface{}, location string) error {
	url := baseURL + key + "&q=" + location + "&format=json"

	out, status, err := api.RequestData(url)

	if err != nil && status != http.StatusOK {
		return err
	}

	// should we use NewDecoder?
	err = json.Unmarshal(out, locations)

	return err
}