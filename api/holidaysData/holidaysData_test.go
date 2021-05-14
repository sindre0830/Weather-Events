package holidaysData

import (
	"net/http"
	"testing"
)

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
