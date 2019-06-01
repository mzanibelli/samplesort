package analyze

import (
	"encoding/json"
	"math"

	"github.com/bugra/kmeans"
)

type dataset interface {
	Features() [][]float64
	Sort([]int)
}

type engine interface {
	Compute([][]float64)
	SDs() []float64
}

type cache interface {
	Fetch(key string, target interface{}, build func() ([]byte, error)) error
}
type Analyze struct {
	data      dataset
	stats     engine
	cache     cache
	size      int
	threshold int
}

func New(data dataset, stats engine, cache cache, size, threshold int) *Analyze {
	return &Analyze{
		data,
		stats,
		cache,
		size,
		threshold,
	}
}

func (a *Analyze) Analyze() error {
	var feats [][]float64
	var result []int
	var err error
	err = a.cache.Fetch("features", &feats,
		func() ([]byte, error) {
			feats = a.data.Features()
			return json.Marshal(feats)
		})
	if err != nil {
		return err
	}
	err = a.cache.Fetch("kmeans", &result,
		func() ([]byte, error) {
			a.stats.Compute(feats)
			dist := a.Distance(a.stats.SDs())
			result, err = kmeans.Kmeans(feats, a.size, dist, a.threshold)
			if err != nil {
				return nil, err
			}
			return json.Marshal(result)
		})
	if err != nil {
		return err
	}
	a.data.Sort(result)
	return nil
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
