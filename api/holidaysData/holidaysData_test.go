package holidaysData

import (
	"main/fun"
	"main/storage"
	"net/http"
	"testing"
)

func TestHandler(t *testing.T) {
	//change directory to root
	newDir, err := fun.GoToRoot()
	if err != nil {
		t.Fatal(err)
	}

	//Mocked firebase
	err = storage.Firebase.Setup()
	if err != nil {
		defer storage.Firebase.Client.Close()
		t.Error(err)
	}

	//store expected data to check against
	testData := map[string]int{
		"Oslo":              http.StatusOK,
		"Bergen":            http.StatusOK,
		"PlaceDoesntExist":  http.StatusNotFound,
		"There/are/slashes": http.StatusNotFound,
		//Add more cases here
	}
	//iterate through map and check each key to expected element
	for test, expectedStatus := range testData {

		_, status, _ := Handler(test)
		//branch if we get an unexpected answer that is not timed out
		if status != expectedStatus && status != http.StatusRequestTimeout {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'. Err: %v. Path: %v", expectedStatus, status, test, err, newDir)
		}
	}
}

func TestGet(t *testing.T) {
	//store expected data to check against
	testData := map[string]int{
		"no":                http.StatusOK,
		"Dont/put/slashes/": http.StatusNotFound,
		"dk":                http.StatusOK,
		"iLiGasghFYG":       http.StatusNotFound,
	}
	//iterate through map and check each key to expected element
	for test, expectedStatus := range testData {
		_, status, _ := get(test)
		//branch if we get an unexpected answer that is not timed out
		if status != expectedStatus && status != http.StatusRequestTimeout {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'.", expectedStatus, status, test)
		}
	}
}
