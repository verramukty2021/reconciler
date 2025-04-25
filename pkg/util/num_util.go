package util

import (
	"math"
)

func RoundToTwoDecimals(num float64) float64 {
	rounded := math.Round(num*100) / 100
	return rounded
}
