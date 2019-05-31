package engine

import "math"

func round(n, precision float64) float64 {
	if precision == 0 {
		precision = math.SmallestNonzeroFloat64
	}
	return math.Round(n/precision) * precision
}

func sum(data []float64) float64 {
	var result float64
	for _, value := range data {
		result += value
	}
	return result
}

func variance(data []float64) float64 {
	if len(data) < 2 {
		return 0
	}
	var K, n, Ex, Ex2 float64
	K = data[0]
	for _, x := range data {
		n++
		Ex += x - K
		Ex2 += (x - K) * (x - K)
	}
	return (Ex2 - (Ex*Ex)/n) / n
}
