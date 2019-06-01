package engine

import (
	"encoding/json"
	"math"
)

type Engine struct {
	stats map[int]*stat
}

func New() *Engine {
	return new(Engine)
}

func (e *Engine) Compute(data [][]float64) {
	for i, features := range data {
		if i == 0 {
			e.stats = make(map[int]*stat, len(features))
		}
		for j, feat := range features {
			if _, ok := e.stats[j]; !ok {
				e.stats[j] = &stat{
					values: make([]float64, len(data), len(data)),
					Min:    +math.MaxFloat64,
					Max:    -math.MaxFloat64,
				}
			}
			e.stats[j].update(i, feat)
		}
	}
}

func (e *Engine) SDs() []float64 {
	res := make([]float64, len(e.stats))
	for i := range res {
		res[i] = e.stats[i].sd()
	}
	return res
}

func (e *Engine) Means() []float64 {
	res := make([]float64, len(e.stats))
	for i := range res {
		res[i] = e.stats[i].mean()
	}
	return res
}

func (e *Engine) String() string {
	json, _ := json.MarshalIndent(e.stats, "", " ")
	return string(json)
}

type stat struct {
	values []float64
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
}

func (s *stat) update(i int, n float64) {
	s.values[i] = n
	s.Min = math.Min(s.Min, n)
	s.Max = math.Max(s.Max, n)
}

func (s *stat) sd() float64 {
	return math.Sqrt(variance(s.values))
}

func (s *stat) mean() float64 {
	return mean(s.values)
}
