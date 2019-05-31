package analyze

import (
	"math"

	"github.com/bugra/kmeans"
)

type Analyze struct {
	data      dataset
	stats     engine
	size      int
	threshold int
}

type dataset interface {
	Features() [][]float64
	Sort([]int)
}

type engine interface {
	Compute([][]float64)
	SDs() []float64
}

func New(data dataset, stats engine, size, threshold int) *Analyze {
	return &Analyze{
		data,
		stats,
		size,
		threshold,
	}
}

func (a *Analyze) Analyze() {
	feats := a.data.Features()
	a.stats.Compute(feats)
	dist := a.Distance(a.stats.SDs())
	res, err := kmeans.Kmeans(feats, a.size, dist, a.threshold)
	if err != nil {
		panic(err)
	}
	a.data.Sort(res)
}

// Distance is an Hamming distance that tolerates an error margin
// testing float values for equality.
func (a *Analyze) Distance(margins []float64) kmeans.DistanceFunction {
	return func(s1, s2 []float64) (float64, error) {
		var res float64 = 0
		for i, margin := range margins {
			if math.Abs(s1[i]-s2[i]) > margin/2 {
				res++
			}
		}
		return res, nil
	}
}
