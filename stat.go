package samplesort

import "math"

type Stat struct {
	values []float64
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
	Mean   float64 `json:"mean"`
}

func (s *Stat) Update(value float64) {
	if s.values == nil {
		s.values = make([]float64, 0)
	}
	s.values = append(s.values, value)
	s.Min = math.Min(s.Min, value)
	s.Max = math.Max(s.Max, value)
	s.Mean = sum(s.values) / float64(len(s.values))
}

func sum(values []float64) float64 {
	var result float64
	for _, value := range values {
		result += value
	}
	return result
}
