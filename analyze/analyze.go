package analyze

import (
	"github.com/bugra/kmeans"
)

type dataset interface {
	Features() [][]float64
	Sort([]int)
}

type engine interface {
	Normalize([][]float64) [][]float64
	Distance(s1, s2 []float64) (float64, error)
}

type cache interface {
	Fetch(key string, target interface{}, build func() ([]byte, error)) error
	Serialize(v interface{}) ([]byte, error)
}

type config interface {
	Size() int
	MaxIterations() int
	Log(vs ...interface{})
}

type Analyze struct {
	data  dataset
	stats engine
	cache cache
	cfg   config
}

func New(data dataset, stats engine, storage cache, cfg config) *Analyze {
	return &Analyze{
		data:  data,
		stats: stats,
		cache: storage,
		cfg:   cfg,
	}
}

func (a *Analyze) Analyze() ([][]float64, error) {
	var rawFeatures [][]float64
	var normalizedFeatures [][]float64
	var result []int
	var err error

	a.cfg.Log("gathering features...")
	err = a.cache.Fetch("features", &rawFeatures,
		func() ([]byte, error) {
			rawFeatures = a.data.Features()
			return a.cache.Serialize(rawFeatures)
		})
	if err != nil {
		return nil, err
	}

	a.cfg.Log("normalizing features...")
	normalizedFeatures = a.stats.Normalize(rawFeatures)

	a.cfg.Log("computing kmeans...")
	err = a.cache.Fetch("kmeans", &result,
		func() ([]byte, error) {
			result, err = kmeans.Kmeans(normalizedFeatures, a.cfg.Size(),
				kmeans.SquaredEuclideanDistance, a.cfg.MaxIterations())
			if err != nil {
				return nil, err
			}
			return a.cache.Serialize(result)
		})
	if err != nil {
		return nil, err
	}

	a.cfg.Log("sorting dataset...")
	a.data.Sort(result)

	return normalizedFeatures, nil
}
