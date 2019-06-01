package engine

import (
	"encoding/json"
	"math"

	"gonum.org/v1/gonum/stat"
)

const (
	ZSCORE_MIN float64 = -2.0
	ZSCORE_MAX float64 = 2.0
)

type Engine struct {
	stats map[int]*featStat
}

func New() *Engine {
	return new(Engine)
}

func (e *Engine) Compute(data [][]float64) {
	if len(e.stats) > 0 {
		return
	}
	for i, features := range data {
		if i == 0 {
			e.stats = make(map[int]*featStat, len(features))
		}
		for j, feat := range features {
			if _, ok := e.stats[j]; !ok {
				e.stats[j] = &featStat{
					values: make([]float64, len(data), len(data)),
				}
			}
			e.stats[j].update(i, feat)
		}
	}
}

func (e *Engine) Normalize(data [][]float64) [][]float64 {
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
		_, std := e.stats[i].meanStdDev()
		if math.Abs(s1[i]-s2[i]) > std/2 {
			res++
		}
	}
	return res, nil
}

func (e *Engine) String() string {
	json, _ := json.MarshalIndent(e.stats, "", " ")
	return string(json)
}

type featStat struct {
	values []float64
	mean   float64
	std    float64
}

func (s *featStat) meanStdDev() (mean, std float64) {
	if s.mean == 0 && s.std == 0 {
		s.mean, s.std = stat.MeanStdDev(s.values, s.weights())
	}
	return s.mean, s.std
}

func (s *featStat) stdScore(x float64) float64 {
	mean, std := s.meanStdDev()
	return stat.StdScore(x, mean, std)
}

func (f *featStat) normalize(n float64) float64 {
	score := f.stdScore(n)
	switch {
	case ZSCORE_MIN > score:
		n = ZSCORE_MIN * n / score
		break
	case score > ZSCORE_MAX:
		n = ZSCORE_MAX * n / score
		break
	}
	return n
}

func (s *featStat) update(i int, n float64) { s.values[i] = n }

func (s *featStat) weights() []float64 {
	weigths := make([]float64, len(s.values))
	for i := range weigths {
		weigths[i] = 1 // so far, we cannot weight features
	}
	return weigths
}
