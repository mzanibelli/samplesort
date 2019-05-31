package engine

import (
	"encoding/json"
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
		e.stats[key] = &stat{
			name:   key,
			values: make([]float64, 0),
			Min:    +math.MaxFloat64,
			Max:    -math.MaxFloat64,
			Mean:   0,
		}
	}
	e.stats[key].update(n)
}

func (e *Engine) Normalize(key string, n float64) float64 {
	s, ok := e.stats[key]
	switch {
	case !ok:
		return 0
	case s.Max == s.Min:
		return 0
	case s.Min == 0 && s.Max == 1:
		return e.round(n)
	default:
		return e.round((n - s.Min) / (s.Max - s.Min))
	}
}

func (e *Engine) String() string {
	json, _ := json.MarshalIndent(e.stats, "", " ")
	return string(json)
}

func (e *Engine) round(n float64) float64 {
	return math.Round(n/e.precision) * e.precision
}

type stat struct {
	name   string
	values []float64
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
	Mean   float64 `json:"mean"`
}

func (s *stat) String() string { return s.name }

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
