package dict

import "testing"

func TestValidGetYrUrl(t *testing.T) {
	data := GetYrURL("45.4", "30.1")

	if data != YR_URL+"?lat="+"45.4"+"&lon="+"30.1" {
		t.Errorf("testHandlerValid failed, expected %v, got %v", "https://api.met.no/weatherapi/locationforecast/2.0/complete?lat=45.4&lon=30.1", data)
	} else {
		t.Logf("testHandlerValid success, expected %v, got %v", "https://api.met.no/weatherapi/locationforecast/2.0/complete?lat=45.4&lon=30.1", data)
	}

}

//The other tests for this package would be very similar to the other test
