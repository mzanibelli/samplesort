package engine

import (
	"math"
)

type Engine struct {
	precision float64
	stats     map[string]*stat
}

func New(precision float64) *Engine {
	return &Engine{precision, make(map[string]*stat, 0)}
}

func (e *Engine) Update(key string, n float64) {
	if _, ok := e.stats[key]; !ok {
		e.stats[key] = &stat{make([]float64, 0), 0, 0, 0}
	}
	e.stats[key].update(n)
}

func (e *Engine) Normalize(key string, n float64) float64 {
	if _, ok := e.stats[key]; !ok {
		return 0
	}
	min := e.stats[key].Min
	max := e.stats[key].Max
	if (max - min) == 0 {
		return 0
	}
	norm := (n - min) / (max - min)
	x := math.Round(norm/e.precision) * e.precision
	return x
}

type stat struct {
	values []float64
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
	Mean   float64 `json:"mean"`
}

func (s *stat) update(n float64) {
	if s.values == nil {
		s.values = make([]float64, 0)
	}
	s.values = append(s.values, n)
	s.Min = math.Min(s.Min, n)
	s.Max = math.Max(s.Max, n)
	s.Mean = sum(s.values) / float64(len(s.values))
}

func sum(values []float64) float64 {
	var result float64
	for _, value := range values {
		result += value
	}
	return result
}
