package fun

import "testing"

func TestLimitDecimals(t *testing.T) {
	testData := map[float64]float64{
		45.57859:         45.58,
		107.564329486328: 107.56,
	}
	//iterate through map and check each key to expected element
	for test, expectedResult := range testData {
		result := LimitDecimals(test)
		//branch if we get an unexpected answer
		if result != expectedResult {
			t.Errorf("Expected '%v' but got '%v'. Tested: '%v'", expectedResult, result, test)
		}
	}
}
