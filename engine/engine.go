package engine

import (
	"encoding/json"
	"math"

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
		res += math.Abs(sampleFeatures[i] - meanOfCluster[i])
	}
	return res / float64(len(sampleFeatures)), nil
}

func (e *Engine) String() string {
	str, _ := json.MarshalIndent(e.stats, "", " ")
	return string(str)
}

// See: https://en.wikipedia.org/wiki/Feature_scaling
// The general method of calculation is to determine the distribution
// mean and standard deviation for each feature. Next we subtract the
// mean from each feature. Then we divide the values (mean is already
// subtracted) of each feature by its standard deviation.
func (e *Engine) Normalize(data [][]float64) func(i, j int, v float64) float64 {
	e.feed(data)
	return func(i, j int, v float64) float64 {
		return (v - e.stats[j].mean) / e.stats[j].std
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
		s.min = math.Min(s.min, v)
		s.max = math.Max(s.max, v)
	}
}

// TODO: how to smartly weight features?
func (s *featStat) weights() []float64 { return nil }
