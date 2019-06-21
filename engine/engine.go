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
	return &Engine{
		cfg: cfg,
	}
}

// TODO: continue investigation.
func (e *Engine) Distance(sampleFeatures, meanOfCluster []float64) (float64, error) {
	var res float64 = 0
	for i := range sampleFeatures {
		res += math.Abs(sampleFeatures[i] - meanOfCluster[i])
	}
	return res, nil
}

func (e *Engine) String() string {
	json, _ := json.MarshalIndent(e.stats, "", " ")
	return string(json)
}

func (e *Engine) Normalize(data [][]float64) [][]float64 {
	e.feed(data)
	res := make([][]float64, len(data), len(data))
	for i := range res {
		row := make([]float64, len(e.stats), len(e.stats))
		for j := range row {
			row[j] = e.stats[j].normalize(data[i][j])
		}
		res[i] = row
	}
	return res
}

func (e *Engine) feed(data [][]float64) {
	for i, features := range data {
		if i == 0 {
			e.stats = make(map[int]*featStat, len(features))
		}
		for j, feat := range features {
			e.update(i, j, len(data), feat)
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
	warm   bool
}

func newFeatStat(size int) *featStat {
	return &featStat{
		values: make([]float64, size, size),
		mean:   0,
		std:    0,
		min:    math.MaxFloat64,
		max:    -math.MaxFloat64,
		warm:   false,
	}
}

func (s *featStat) setMeanStd() {
	s.mean, s.std = stat.MeanStdDev(s.values, s.weights())
}

func (s *featStat) setMinMax(zscore float64) {
	for _, v := range s.values {
		if math.Abs(stat.StdScore(v, s.mean, s.std)) > zscore {
			continue
		}
		s.min = math.Min(s.min, v)
		s.max = math.Max(s.max, v)
	}
}

func (s *featStat) normalize(n float64) float64 {
	norm := (n - s.min) / (s.max - s.min)
	switch {
	case norm >= 1:
		return 1
	case norm <= 0:
		return 0
	default:
		return norm
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
