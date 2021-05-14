package eventData

import (
	"net/http"
	"testing"
)

func TestReq(t *testing.T) {

	//store expected data to check against
	testData := map[string]int{
		"https://https://app.ticketmaster.com/discovery/v2/events/Z698xZbpZ17a4oM.json?apikey=ySyIqc6FFKgUIIgzKB5LAOQeGUeU1mot": http.StatusOK,
		"https://notok/wontwork": http.StatusInternalServerError,
		"https://https://app.ticketmaster.com/discovery/v2/events/vvG1YZ4VLloAHj.json?apikey=ySyIqc6FFKgUIIgzKB5LAOQeGUeU1mot": http.StatusOK,
		"https://jsonapi.org/examples/": http.StatusInternalServerError,
	}
	//iterate through map and check each key to expected element
	for test, expectedStatus := range testData {
		var data EventInformation
		status, _ := data.req(test)
		//branch if we get an unexpected answer that is not timed out
		if status != expectedStatus && status != http.StatusRequestTimeout {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'.", expectedStatus, status, test)
		}
	}

}
