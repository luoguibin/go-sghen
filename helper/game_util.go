package helper

import(
	"math"
)

func GClientDistance(x0 int, y0 int, x1 int, y1 int) float64 {
	xV := math.Pow(float64(x0 - x1), 2.0)
	yV := math.Pow(float64(y0 - y1), 2.0)
	return math.Sqrt(xV + yV)
}