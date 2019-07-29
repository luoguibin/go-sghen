package helper

import (
	"math"
)

func GClientDistance(x0 float64, y0 float64, x1 float64, y1 float64) float64 {
	xV := math.Pow(x0-x1, 2.0)
	yV := math.Pow(y0-y1, 2.0)
	return math.Sqrt(xV + yV)
}

func GSliceRemove(s []int, index int) []int {
	if index < 0 || index >= len(s) {
		return s
	}
	if index == len(s)-1 {
		return s[0:index]
	} else {
		return append(s[0:index], s[index+1:]...)
	}
}
