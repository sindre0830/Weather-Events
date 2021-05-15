package weatherCompare

import (
	"encoding/json"
	"main/dict"
	"main/fun"
	"main/storage"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	//change directory to root
	_, err := fun.GoToRoot()
	if err != nil {
		t.Fatal(err)
	}

	//Mocked firebase
	err = storage.Firebase.Setup()
	if err != nil {
		defer storage.Firebase.Client.Close()
		t.Error(err)
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", dict.WEATHERCOMPARE_PATH+"oslo"+"/bergen;stavanger", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(MethodHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var input WeatherCompare
	err = json.NewDecoder(rr.Body).Decode(&input)
	if err != nil {
		t.Fatal(err)
	}

	// Check the response body is what we expect
	if input.Latitude < 0 || input.Latitude > 180 {
		t.Error("Error in formatting of responce struct")
	}
	//Add more cases here

}

func TestGet(t *testing.T) {
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

	//Checks
	var data WeatherCompare
	var array = []locationInfo{
		locationInfo{10.74, 59.91, "oslo"},
		locationInfo{5.71, 59.1, "stavanger"}}
	status, err := data.get(35.5, 23.6, array, time.Now().AddDate(0, 0, 1).Format("2006-01-02")) //should return status ok
	//branch if we get an unexpected answer that is not timed out
	if status != http.StatusOK && status != http.StatusRequestTimeout {
		t.Errorf("Expected '%v' but got '%v'. Tested: '%v'. Err: %v. Path: %v", http.StatusOK, status, "35.5, 23.6, \"current date + 1 day\"", err, newDir)
	}

	array = []locationInfo{
		locationInfo{5.70, 59.11, "stavanger"},
		locationInfo{5.33, 60.39, "bergen"}}
	status, err = data.get(35.5, 23.6, array, time.Now().AddDate(0, 0, 1).Format("2006-01-02")) //should return status ok
	//branch if we get an unexpected answer that is not timed out
	if status != http.StatusOK && status != http.StatusRequestTimeout {
		t.Errorf("Expected '%v' but got '%v'. Tested: '%v'. Err: %v. Path: %v", http.StatusOK, status, "35.5, 23.6, \"current date + 1 day\"", err, newDir)
	}
}
