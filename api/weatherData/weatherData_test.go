package weatherData

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
	testData := map[[2]string]int{
		{"9.4", "23.4"}: http.StatusOK,
		//Add more cases here
	}
	//iterate through map and check each key to expected element
	for test, expectedStatus := range testData {
		var data WeatherData
		status, err := data.Handler(test[0], test[1])
		//branch if we get an unexpected answer that is not timed out
		if status != expectedStatus && status != http.StatusRequestTimeout {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'. Err: %v. Path: %v", expectedStatus, status, test, err, newDir)
		}
	}
}
