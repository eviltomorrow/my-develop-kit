package zmath

import "math"

// IRR 计算(错误)
func IRR(total float64, per float64, n int) (float64, float64) {
	if total <= 0 {
		return 0, 0
	}
	var rate float64
	var npv = total + per*float64(n)/math.Pow(1+rate, float64(n))
	return IRR(npv, per, n+1)
}
