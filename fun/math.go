package fun

import "math"

// LimitDecimals limits decimals to two.
func LimitDecimals(number float64) float64 {
	return math.Round(number * 100) / 100
}
