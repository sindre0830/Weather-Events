package weatherEvent

import (
	"main/dict"
	"main/fun"
	"main/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAll(t *testing.T) {
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
	req, err := http.NewRequest("GET", dict.WEATHEREVENT_HOOK_PATH, nil)
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

func TestCheckDate(t *testing.T) {

	var data WeatherEvent
	if data.checkDate() { //branch if we get an unexpected answer
		t.Errorf("Expected '%v' but got '%v'. Tested: '%v'", "true", data.checkDate(), "data.checkDate()")
	}

}
