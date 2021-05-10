package geocoords

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"main/api"
	"main/db"
	"main/debug"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

/**
*	LocationCoords
*	Holds our latitude and longitude data for one location
**/
type LocationCoords struct {
	Address    string  `json:"address"`
	Importance float64 `json:"importance"`
	Latitude   float64 `json:"lat"`
	Longitude  float64 `json:"lon"`
}

var baseURL = "https://us1.locationiq.com/v1/search.php?key="
var key = "pk.d8a67c78822d16869c7a3e8f6d7617af"

// This loads in local database from file of most important locations
var LocalCoords = make(map[string]LocationCoords)

/**
*	CoordHandler
*	Accepts a place name, and gets the geological coordinages associated with that place.
*
*	@param	id			-	String containing the name of the location we want the coordinates of
*	@return	int, error)	-	Tuple with an http status code and an error interface. If everything works ok, the error is nil
*
*	@see	getCoords
*	@see	getLocations
**/
func (locationCoords *LocationCoords) Handler(id string) (int, error) {
	// We read our local DB, if one exists, into LocalCoords map.
	var file []byte
	if _, err := os.Stat("GeoCoords.json"); err == nil {
		file, err = ioutil.ReadFile("GeoCoords.json")
		err = json.Unmarshal([]byte(file), &LocalCoords)
	}

	// Check LocalCoords if location data for this location exists
	localData, found := LocalCoords[id]

	if found {
		// If the data was on file, we set it here and return!
		locationCoords.Address = localData.Address
		locationCoords.Importance = localData.Importance
		locationCoords.Latitude = localData.Latitude
		locationCoords.Longitude = localData.Longitude

		return http.StatusOK, nil
	}

	// If not, we check in firestoreDB if location data for this location exists
	data, exist, err := db.DB.Get("GeoCoords", id)
	if err != nil && exist {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"GeoCoords.HandlerCoords() -> Database.get() -> Trying to get data",
			err.Error(),
			"Unknown",
		)
		return http.StatusInternalServerError, err
	}

	// We check whether data on firestore is deprecated or not.
	// For locations that are not countries/capitals, we don't want to keep our data more than 3 hours.
	// Data saved in local files should be kept indefinitely, so we don't check it.
	withinTimeframe, err := db.CheckDate(data.Time, 3)
	if exist && withinTimeframe {
		if err != nil {
			return http.StatusInternalServerError, err
		}

		err = readData(locationCoords, data.Container)

		if err != nil {
			return http.StatusInternalServerError, err
		}

		return http.StatusOK, nil
	}

	// If the location is not stored in firestore OR locally, We get the data from the locationiq api
	var locations []map[string]interface{}
	status, err := getLocations(&locations, id)

	if err != nil {
		return status, err
	}
	// We get lat and lon from the first json object in our locations array
	err = getCoords(locationCoords, locations[0])
	if err != nil {
		return http.StatusInternalServerError, err
	}
	// We store our fresh data
	if locationCoords.Importance > 0.7 {
		// Save locally if it's an important place
		LocalCoords[id] = *locationCoords
		file, err := json.MarshalIndent(LocalCoords, "", " ")

		err = ioutil.WriteFile("GeoCoords.json", file, 0644)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	} else {
		// If not important, we send the data to firestore
		var data db.Data
		data.Time = time.Now().String()
		data.Container = locationCoords
		_, err = db.DB.Add("GeoCoords", id, data)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}
	return http.StatusOK, nil
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

	coords.Address = location["display_name"].(string)
	latitude, err := strconv.ParseFloat(location["lat"].(string), 64)
	coords.Latitude = math.Round(latitude*100) / 100
	longitude, err := strconv.ParseFloat(location["lon"].(string), 64)
	coords.Longitude = math.Round(longitude*100) / 100
	importance := location["importance"].(float64)
	coords.Importance = math.Round(importance*100) / 100

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
		return errors.New("getting data from database: Can't find expected field Latitude")
	}
	if field, ok := m["Longitude"].(float64); ok {
		coords.Longitude = field
	} else {
		return errors.New("getting data from database: Can't find expected field Longitude")
	}
	if field, ok := m["Importance"].(float64); ok {
		coords.Importance = field
	} else {
		return errors.New("getting data from database: Can't find expected field Importance")
	}
	if field, ok := m["Address"].(string); ok {
		coords.Address = field
	} else {
		return errors.New("getting data from database: Can't find expected field Address")
	}
	return nil
}
