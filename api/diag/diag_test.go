package diag

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

//Code from https://www.digitalflapjack.com/blog/2018/4/10/better-testing-for-golang-http-handlers
func TestMethodHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/diag", nil)
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

	var input DiagStatuses
	err = json.NewDecoder(rr.Body).Decode(&input)
	if err != nil {
		t.Fatal(err)
	}

	// Check the response body is what we expect. //NB
	if input.Restcountries < 100 || input.Restcountries > 600 {
		t.Error("Restcountries did not return a status 100 < restcountries < 600")
	}
	if input.LocationIq < 100 || input.LocationIq > 600 {
		t.Error("LocationIq did not return a status 100 < LocationIq < 600")
	}
	if input.PublicHolidays < 100 || input.PublicHolidays > 600 {
		t.Error("PublicHolidays did not return a status 100 < PublicHolidays < 600")
	}
	if input.TicketMaster < 100 || input.TicketMaster > 600 {
		t.Error("TicketMaster did not return a status 100 < TicketMaster < 600")
	}
	if input.Weatherapi < 100 || input.Weatherapi > 600 {
		t.Error("Weatherapi did not return a status 100 < Weatherapi < 600")
	}
	if input.RegisteredWebhooks < 0 || input.RegisteredWebhooks > 10000 {
		t.Error("RegisteredWebhooks did not return a number 0 < RegisteredWebhooks < 10000")
	}
	if input.Uptime < 0 {
		t.Error("Uptime is negative")
	}
	if string([]rune(input.Version)[0]) != "v" {
		t.Error("Version is not written correctly")
	}

}

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
