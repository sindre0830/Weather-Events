package geoCoords

import (
	"main/storage"
	"net/http"
	"os"
	"testing"
)

func TestHandler(t *testing.T) {

	//Change directory
	os.Chdir("./../../")
	newDir, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	//Mocked firebase
	err = storage.Firebase.Setup()
	if err != nil {
		defer storage.Firebase.Client.Close()
		t.Error(err)
	}

	//Store expected data to check against
	testData := map[string]int{
		"oslo":              http.StatusOK,
		"Dont/put/slashes/": http.StatusNotFound,
		"Oslo":              http.StatusOK,
		"iLiGasghFYG":       http.StatusNotFound,
		//Add more cases here
	}
	//Iterate through map and check each key to expected element
	for test, expectedStatus := range testData {
		var data LocationCoords
		status, err := data.Handler(test)
		//branch if we get an unexpected answer that is not timed out
		if status != expectedStatus && status != http.StatusRequestTimeout {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'. Err: %v. Path: %v", expectedStatus, status, test, err, newDir)
		}
	}
}

func TestGetCoords(t *testing.T) {
	//store expected data to check against
	testData := map[string]int{
		"Oslo": http.StatusOK,
		//Add more cases here
	}

	//iterate through map and check each key to expected element
	for test, expectedStatus := range testData {

		var locations []map[string]interface{}
		status, err := getLocations(&locations, test)
		if err != nil {
			t.Errorf("testHandlerValid failed, expected %v, got %v", "nil", err)
		} else {
			t.Logf("testHandlerValid success, expected %v, got %v", "nil", err)
		}

		var data LocationCoords
		err = getCoords(&data, locations[0])

		if err != nil {
			t.Errorf("testHandlerValid failed, expected %v, got %v", "nil", err)
		} else {
			t.Logf("testHandlerValid success, expected %v, got %v", "nil", err)
		}
		if status != expectedStatus && status != http.StatusRequestTimeout {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'.", expectedStatus, status, test)
		}
	}

}

func TestGetLocations(t *testing.T) {
	//store expected data to check against
	testData := map[string]int{
		"oslo":              http.StatusOK,
		"Dont/put/slashes/": http.StatusNotFound,
		"Oslo":              http.StatusOK,
		"iLiGasghFYG":       http.StatusNotFound,
		//Add more cases here
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
