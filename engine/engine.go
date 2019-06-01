package engine

import (
	"encoding/json"
	"math"

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
	return res
}

func (e *Engine) Distance(s1, s2 []float64) (float64, error) {
	var res float64 = 0
	for i := range e.stats {
		if math.Abs(s1[i]-s2[i]) > e.stats[i].std/2 {
			res++
		}
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
		e.stats[j] = new(featStat)
		e.stats[j].values = make([]float64, size, size)
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
	score := math.Abs(stat.StdScore(n, s.mean, s.std))
	if score > ZSCORE_MAX {
		n = ZSCORE_MAX * n / score
	}
	return n
}

// so far, we cannot weight features
func (s *featStat) weights() []float64 {
	weigths := make([]float64, len(s.values))
	for i := range weigths {
		weigths[i] = 1
	}
	return weigths
}
