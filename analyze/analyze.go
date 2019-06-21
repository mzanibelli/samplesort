package analyze

import (
	"log"

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

type Analyze struct {
	data      dataset
	stats     engine
	cache     cache
	size      int
	threshold int
	loggers   []*log.Logger
}

func New(
	data dataset,
	stats engine,
	cache cache,
	size, threshold int,
	loggers ...*log.Logger,
) *Analyze {
	a := new(Analyze)
	a.data = data
	a.stats = stats
	a.cache = cache
	a.size = size
	a.threshold = threshold
	a.loggers = loggers
	return a
}

func (a *Analyze) Analyze() ([][]float64, error) {
	var rawFeatures [][]float64
	var normalizedFeatures [][]float64
	var result []int
	var err error

	a.log("gathering features...")
	err = a.cache.Fetch("features", &rawFeatures,
		func() ([]byte, error) {
			rawFeatures = a.data.Features()
			return a.cache.Serialize(rawFeatures)
		})
	if err != nil {
		return nil, err
	}

	a.log("normalizing features...")
	normalizedFeatures = a.stats.Normalize(rawFeatures)

	a.log("computing kmeans...")
	err = a.cache.Fetch("kmeans", &result,
		func() ([]byte, error) {
			result, err = kmeans.Kmeans(normalizedFeatures, a.size,
				kmeans.SquaredEuclideanDistance, a.threshold)
			if err != nil {
				return nil, err
			}
			return a.cache.Serialize(result)
		})
	if err != nil {
		return nil, err
	}

	a.log("sorting dataset...")
	a.data.Sort(result)

	return normalizedFeatures, nil
}

func (a *Analyze) log(vs ...interface{}) {
	for _, l := range a.loggers {
		l.Println(vs...)
	}
}
