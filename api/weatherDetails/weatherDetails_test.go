package weatherDetails

import (
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
	req, err := http.NewRequest("GET", dict.WEATHERDETAILS_PATH+"oslo", nil)
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
	var data WeatherDetails
	status, err := data.get(35.5, 23.6, time.Now().AddDate(0, 0, 1).Format("2006-01-02")) //should return status ok
	//branch if we get an unexpected answer that is not timed out
	if status != http.StatusOK && status != http.StatusRequestTimeout {
		t.Errorf("Expected '%v' but got '%v'. Tested: '%v'. Err: %v. Path: %v", http.StatusOK, status, "35.5, 23.6, \"current date + 1 day\"", err, newDir)
	}

	status, err = data.get(600, 23.6, "2021-05-17") //returns badrequest because lat is too big
	//branch if we get an unexpected answer that is not timed out
	if status != http.StatusBadRequest && status != http.StatusRequestTimeout {
		t.Errorf("Expected '%v' but got '%v'. Tested: '%v'. Err: %v. Path: %v", http.StatusBadRequest, status, "35.5, 23.6, \"2021-04-04\"", err, newDir)
	}

	status, err = data.get(600, 23.6, "2021-04-04") //returns badrequest because date has passed
	//branch if we get an unexpected answer that is not timed out
	if status != http.StatusBadRequest && status != http.StatusRequestTimeout {
		t.Errorf("Expected '%v' but got '%v'. Tested: '%v'. Err: %v. Path: %v", http.StatusBadRequest, status, "35.5, 23.6, \"2021-04-04\"", err, newDir)
	}

}
