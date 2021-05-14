package diag

import (
	"net/http"
	"testing"
)

func TestGetStatusOf(t *testing.T) {
	testData := map[string]int{
		"https://restcountries.eu/rest/v2/all":                                          http.StatusOK,
		"https://notValid/rest/v2/all":                                                  http.StatusInternalServerError,
		"https://api.met.no/weatherapi/locationforecast/2.0/complete?lat=30.0&lon=30.0": http.StatusForbidden,
		"https://AnnotherInvalid/one/":                                                  http.StatusInternalServerError,
	}

	for test, expectedStatus := range testData {
		status, _ := getStatusOf(test)
		//branch if we get an unexpected answer that is not timed out
		if status != expectedStatus && status != http.StatusRequestTimeout {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'.", expectedStatus, status, test)
		}
	}

}
