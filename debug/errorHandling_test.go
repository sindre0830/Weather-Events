package debug

import "testing"

func TestUpdate(t *testing.T) {
	var ErrorMessage Debug
	ErrorMessage.Update(200, "Handler", "error", "reason")

	if ErrorMessage.StatusCode != 200 {
		t.Errorf("testHandlerValid failed, expected %v, got %v", 200, ErrorMessage.StatusCode)
	} else {
		t.Logf("testHandlerValid success, expected %v, got %v", 200, ErrorMessage.StatusCode)
	}

	if ErrorMessage.Location != "Handler" {
		t.Errorf("testHandlerValid failed, expected %v, got %v", "Handler", ErrorMessage.Location)
	} else {
		t.Logf("testHandlerValid success, expected %v, got %v", "Handler", ErrorMessage.Location)
	}

	if ErrorMessage.RawError != "error" {
		t.Errorf("testHandlerValid failed, expected %v, got %v", "error", ErrorMessage.RawError)
	} else {
		t.Logf("testHandlerValid success, expected %v, got %v", "error", ErrorMessage.RawError)
	}

	if ErrorMessage.PossibleReason != "Unknown" {
		t.Errorf("testHandlerValid failed, expected %v, got %v", "Unknown", ErrorMessage.PossibleReason)
	} else {
		t.Logf("testHandlerValid success, expected %v, got %v", "Unknown", ErrorMessage.PossibleReason)
	}

}
