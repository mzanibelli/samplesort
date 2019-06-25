package engine

import (
	"encoding/json"
	"math"

	"github.com/bugra/kmeans"
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
	return kmeans.HammingDistance(sampleFeatures, meanOfCluster)
}

func (e *Engine) String() string {
	json, _ := json.MarshalIndent(e.stats, "", " ")
	return string(json)
}

func (e *Engine) Normalize(data [][]float64) func(i, j int, v float64) float64 {
	e.feed(data)
	return func(i, j int, v float64) float64 {
		switch {
		case e.stats[j].min < 0:
			return math.Max(0, v+e.stats[j].min)
		case e.stats[j].min > 0:
			return math.Max(0, v-e.stats[j].min)
		default:
			return math.Max(0, v)
		}
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

// TODO: how to smartly weight features?
func (s *featStat) weights() []float64 { return nil }
