package storage

import (
	"testing"
	"time"
)

//Should Implement more efficient method

// func TestCheckDate(t *testing.T) {
// 	testData := map[string]int{
// 		"13 May 21 14:25 CEST": 3,
// 		"13 May 21":            6,
// 	}

// }

func TestValidCheckDate(t *testing.T) {
	status, err := CheckDate("13 May 21 14:25 CEST", 3)

	if status != true {
		t.Errorf("testHandlerValid failed, expected %v, got %v", 200, status)
	} else {
		t.Logf("testHandlerValid success, expected %v, got %v", 200, status)
	}

	if err != nil {
		t.Errorf("testHandlerValid failed, expected %v, got %v", nil, err)
	} else {
		t.Logf("testHandlerValid success, expected %v, got %v", nil, err)
	}
}

func TestInValidCheckDate(t *testing.T) {
	status, err := CheckDate("13 May 21", 6)

	if status != true {
		t.Errorf("testHandlerValid failed, expected %v, got %v", 200, status)
	} else {
		t.Logf("testHandlerValid success, expected %v, got %v", 200, status)
	}

	if err == nil {
		t.Errorf("testHandlerValid failed, expected %v, got %v", "parsing time \"13 May 21\" as \"02 Jan 06 15:04 MST\": cannot parse \"\" as \"15\"", err)
	} else {
		t.Logf("testHandlerValid success, expected %v, got %v", nil, err)
	}
}

func TestValidCheckIfDateOfEventPassed(t *testing.T) {
	status := CheckIfDateOfEventPassed(time.Date(2021, time.Month(2), 21, 1, 10, 30, 0, time.UTC))

	if status != true {
		t.Errorf("testHandlerValid failed, expected %v, got %v", 200, status)
	} else {
		t.Logf("testHandlerValid success, expected %v, got %v", 200, status)
	}

}
