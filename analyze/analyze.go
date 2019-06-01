package analyze

import (
	"encoding/json"
	"log"

	"github.com/bugra/kmeans"
)

type dataset interface {
	Features() [][]float64
	Sort([]int)
}

type engine interface {
	Compute([][]float64)
	Normalize([][]float64) [][]float64
	Distance(s1, s2 []float64) (float64, error)
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
	var rawFeatures [][]float64
	var normalizedFeatures [][]float64
	var result []int
	var err error

	err = a.cache.Fetch("features", &rawFeatures,
		func() ([]byte, error) {
			log.Println("gathering features...")
			rawFeatures = a.data.Features()
			return json.Marshal(rawFeatures)
		})
	if err != nil {
		return err
	}

	err = a.cache.Fetch("normalized", &normalizedFeatures,
		func() ([]byte, error) {
			log.Println("computing features...")
			a.stats.Compute(rawFeatures)
			log.Println("normalizing features...")
			normalizedFeatures = a.stats.Normalize(rawFeatures)
			return json.Marshal(normalizedFeatures)
		})
	if err != nil {
		return err
	}

	err = a.cache.Fetch("kmeans", &result,
		func() ([]byte, error) {
			log.Println("computing features...")
			a.stats.Compute(rawFeatures)
			log.Println("computing kmeans...")
			result, err = kmeans.Kmeans(normalizedFeatures, a.size,
				a.stats.Distance, a.threshold)
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
