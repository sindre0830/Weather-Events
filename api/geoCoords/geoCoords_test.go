package geoCoords

import (
	"net/http"
	"testing"
)

func TestGetLocations(t *testing.T) {
	//store expected data to check against
	testData := map[string]int{
		"oslo":              http.StatusOK,
		"Dont/put/slashes/": http.StatusNotFound,
		"Oslo":              http.StatusOK,
		"iLiGasghFYG":       http.StatusNotFound,
	}
	//iterate through map and check each key to expected element
	for test, expectedStatus := range testData {
		var locations []map[string]interface{}
		status, _ := getLocations(&locations, test)
		//branch if we get an unexpected answer that is not timed out
		if status != expectedStatus && status != http.StatusRequestTimeout {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'.", expectedStatus, status, test)
		}
	}

}
