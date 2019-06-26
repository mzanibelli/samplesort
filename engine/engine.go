package engine

import (
	"encoding/json"
	"math"
)

type Engine struct {
	stats map[int]*featStat
}

func New() *Engine {
	return &Engine{}
}

func (e *Engine) String() string {
	str, _ := json.MarshalIndent(e.stats, "", " ")
	return string(str)
}

func (e *Engine) Distance(sampleFeatures, meanOfCluster []float64) (float64, error) {
	var res float64 = 0
	for i := range sampleFeatures {
		if math.IsNaN(sampleFeatures[i]) {
			continue
		}
		if math.IsNaN(meanOfCluster[i]) {
			continue
		}
		res += math.Abs(sampleFeatures[i] - meanOfCluster[i])
	}
	return res / float64(len(sampleFeatures)), nil
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
	for i := range data {
		for j := range data[i] {
			e.update(i, j, size, data[i][j])
		}
	}
}

func (e *Engine) update(i, j, size int, feat float64) {
	if _, ok := e.stats[j]; !ok {
		e.stats[j] = newFeatStat(size)
	}
	e.stats[j].values[i] = feat
	e.stats[j].setMeanStd()
}
