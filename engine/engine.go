package engine

import (
	"encoding/json"
	"math"
	"sort"

	"gonum.org/v1/gonum/stat"
)

const (
	ZSCORE_MAX float64 = 0.4
)

type Engine struct {
	stats map[int]*featStat
}

func New() *Engine {
	return new(Engine)
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
	e.feed(res)
	return res
}

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
		e.stats[j] = &featStat{
			values: make([]float64, size, size),
		}
	}
	e.stats[j].values[i] = feat
	e.stats[j].mean, e.stats[j].std = stat.MeanStdDev(
		e.stats[j].values,
		e.stats[j].weights(),
	)
}

type featStat struct {
	values []float64
	mean   float64
	std    float64
}

func (s *featStat) normalize(n float64) float64 {
	var min float64 = math.MaxFloat64
	var max float64 = -math.MaxFloat64
	for _, v := range s.values {
		if math.Abs(stat.StdScore(v, s.mean, s.std)) > ZSCORE_MAX {
			continue
		}
		min = math.Min(min, v)
		max = math.Max(max, v)
	}
	norm := (n - min) / (max - min)
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

// so far, we cannot weight features
func (s *featStat) weights() []float64 { return nil }
