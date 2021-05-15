package countryData

import (
	"net/http"
	"os"
	"testing"
)

func TestHandler(t *testing.T) {
	//Change directory code from https://golangbyexample.com/change-current-working-directory-go/#:~:text=Menu-,cd%20command%20in%20Go%20or,working%20directory%20in%20Go%20(Golang)&text=os.,similar%20to%20the%20cd%20command.
	os.Chdir("./../../")
	newDir, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	//store expected data to check against
	testData := map[string]int{
		"Norway": http.StatusOK,
		//Add more cases here
	}
	//iterate through map and check each key to expected element
	for test, expectedStatus := range testData {
		var data Information
		status, err, _ := data.Handler(test)
		//branch if we get an unexpected answer that is not timed out
		if status != expectedStatus && status != http.StatusRequestTimeout {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'. Err: %v. Path: %v", expectedStatus, status, test, err, newDir)
		}
	}

}

func TestReq(t *testing.T) {

	//store expected data to check against
	testData := map[string]int{
		"https://restcountries.eu/rest/v2/all":        http.StatusOK,
		"https://notok/wontwork":                      http.StatusInternalServerError,
		"https://restcountries.eu/rest/v2/name/eesti": http.StatusOK,
		"https://jsonapi.org/examples/":               http.StatusInternalServerError,
		//Add more cases here
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

func TestOneCountry(t *testing.T) {
	//Change directory
	os.Chdir("./../../")
	newDir, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	//store expected data to check against
	testData := map[string]int{
		"Norway": http.StatusOK,
		//Add more cases here
	}
	//iterate through map and check each key to expected element
	for test, expectedStatus := range testData {
		var data Information
		status, err, _ := data.oneCountry("Norway")
		//branch if we get an unexpected answer that is not timed out
		if status != expectedStatus && status != http.StatusRequestTimeout {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'. Err: %v. Path: %v", expectedStatus, status, test, err, newDir)
		}
	}
}

func TestAllCountry(t *testing.T) {
	//Change directory
	os.Chdir("./../../")
	newDir, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	//store expected data to check against
	testData := [1]int{http.StatusOK} //Add more cases here

	//iterate through map and check each key to expected element
	for test := range testData {
		var data Information
		status, err := data.allCountries()
		//branch if we get an unexpected answer that is not timed out
		if status != testData[test] && status != http.StatusRequestTimeout {
			t.Errorf("Got '%v'. Tested: '%v'. Err: %v. Path: %v", status, test, err, newDir)
		}
	}
}
