package fun

import (
	"math"
)

/**
* LimitDecimals
* Truncates a float/double to two decimal numbers.
**/
func LimitDecimals(number float64) float64 {
	return math.Round(number * 100) / 100
}