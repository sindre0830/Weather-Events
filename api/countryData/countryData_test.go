package countryData

import (
	"net/http"
	"testing"
)

func TestReq(t *testing.T) {

	//store expected data to check against
	testData := map[string]int{
		"https://restcountries.eu/rest/v2/all":        http.StatusOK,
		"https://notok/wontwork":                      http.StatusInternalServerError,
		"https://restcountries.eu/rest/v2/name/eesti": http.StatusOK,
		"https://jsonapi.org/examples/":               http.StatusInternalServerError,
	}
	//iterate through map and check each key to expected element
	for test, expectedStatus := range testData {
		var request Information
		status, _ := request.req(test)
		//branch if we get an unexpected answer that is not timed out
		if status != expectedStatus && status != http.StatusRequestTimeout {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'.", expectedStatus, status, test)
		}
	}

}
