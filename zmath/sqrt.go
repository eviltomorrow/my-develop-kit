package zmath

import "math"

// Sqrt sqrt
func Sqrt(x float64) float64 {
	var z = float64(1)
	var tmp = float64(0)

	for math.Abs(tmp-z) > 0.0000000001 {
		tmp = z
		z = (z + x/z) / 2
	}
	return z
}
