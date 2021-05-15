package eventData

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

	//store expected data to check against
	testData := map[string]int{
		"vvG1YZ4VLloAHj": http.StatusOK,
		//Add more cases here
	}
	//iterate through map and check each key to expected element
	for test, expectedStatus := range testData {
		var data EventInformation
		status, err := data.Handler(test)
		//branch if we get an unexpected answer that is not timed out
		if status != expectedStatus && status != http.StatusRequestTimeout {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'. Err: %v. Path: %v", expectedStatus, status, test, err, newDir)
		}
	}
}

func TestGet(t *testing.T) {

	//store expected data to check against
	testData := map[string]int{
		"https://notok/wontwork": http.StatusInternalServerError,
		//Add more cases here
	}
	//iterate through map and check each key to expected element
	for test, expectedStatus := range testData {
		var data EventInformation
		status, _ := data.get(test)
		//branch if we get an unexpected answer that is not timed out
		if status != expectedStatus && status != http.StatusRequestTimeout {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'.", expectedStatus, status, test)
		}
	}

}

func TestReq(t *testing.T) {

	//store expected data to check against
	testData := map[string]int{
		"https://https://app.ticketmaster.com/discovery/v2/events/Z698xZbpZ17a4oM.json?apikey=ySyIqc6FFKgUIIgzKB5LAOQeGUeU1mot": http.StatusOK,
		"https://notok/wontwork": http.StatusInternalServerError,
		"https://https://app.ticketmaster.com/discovery/v2/events/vvG1YZ4VLloAHj.json?apikey=ySyIqc6FFKgUIIgzKB5LAOQeGUeU1mot": http.StatusOK,
		"https://jsonapi.org/examples/": http.StatusInternalServerError,
		//Add more cases here
	}
	//iterate through map and check each key to expected element
	for test, expectedStatus := range testData {
		var data Ticketmaster
		status, _ := data.req(test)
		//branch if we get an unexpected answer that is not timed out
		if status != expectedStatus && status != http.StatusRequestTimeout {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'.", expectedStatus, status, test)
		}
	}

}
