package weatherEvent

import (
	"testing"
)

func TestCheckDate(t *testing.T) {

	var data WeatherEvent
	if data.checkDate() { //branch if we get an unexpected answer
		t.Errorf("Expected '%v' but got '%v'. Tested: '%v'", "true", data.checkDate(), "data.checkDate()")
	}

}
