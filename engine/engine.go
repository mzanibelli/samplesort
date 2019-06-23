package engine

import (
	"encoding/json"
	"math"
	"sort"

	"gonum.org/v1/gonum/stat"
)

type config interface {
	MaxZScore() float64
}

type Engine struct {
	stats map[int]*featStat
	cfg   config
}

func New(cfg config) *Engine {
	return &Engine{cfg: cfg}
}

func (e *Engine) Distance(sampleFeatures, meanOfCluster []float64) (float64, error) {
	var res float64 = 0
	for i := range sampleFeatures {
		difference := math.Abs(sampleFeatures[i] - meanOfCluster[i])
		threshold := e.stats[i].std / 2
		if difference > threshold {
			res += 1
		}
	}
	return res, nil
}

func (e *Engine) String() string {
	json, _ := json.MarshalIndent(e.stats, "", " ")
	return string(json)
}

// TODO: make quantile configurable.
func (e *Engine) Normalize(data [][]float64) func(i, j int, v float64) float64 {
	e.feed(data)
	return func(i, j int, v float64) float64 {
		return math.Min(e.stats[j].quantile(0.75), v)
	}
}

func (e *Engine) feed(data [][]float64) {
	size := len(data)
	if size == 0 {
		return
	}
	e.stats = make(map[int]*featStat, len(data[0]))
	for i, features := range data {
		for j, feat := range features {
			e.update(i, j, size, feat)
		}
	}
}

func (e *Engine) update(i, j, size int, feat float64) {
	if _, ok := e.stats[j]; !ok {
		e.stats[j] = newFeatStat(size)
	}
	e.stats[j].values[i] = feat
	e.stats[j].setMeanStd()
	e.stats[j].setMinMax(e.cfg.MaxZScore())
}

type featStat struct {
	values []float64
	mean   float64
	std    float64
	min    float64
	max    float64
}

func newFeatStat(size int) *featStat {
	return &featStat{
		values: make([]float64, size, size),
		mean:   0,
		std:    0,
		min:    math.MaxFloat64,
		max:    -math.MaxFloat64,
	}
}

func (s *featStat) setMeanStd() {
	s.mean, s.std = stat.MeanStdDev(s.values, s.weights())
}

func (s *featStat) setMinMax(threshold float64) {
	for _, v := range s.values {
		score := math.Abs(stat.StdScore(v, s.mean, s.std))
		if score > threshold {
			continue
		}
		s.min = math.Min(s.min, v)
		s.max = math.Max(s.max, v)
	}
}

func (s *featStat) quantile(n float64) float64 {
	tmp := make([]float64, len(s.values))
	copy(tmp, s.values)
	sort.Float64s(tmp)
	return stat.Quantile(n, stat.Empirical, tmp, s.weights())
}

// TODO: how to smartly weight features?
func (s *featStat) weights() []float64 { return nil }
