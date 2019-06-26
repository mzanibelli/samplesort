package analyze

import "github.com/bugra/kmeans"

type dataset interface {
	Features() [][]float64
	Sort([]int)
}

type engine interface {
	Normalize([][]float64) func(i, j int, v float64) float64
	Distance(s1, s2 []float64) (float64, error)
}

type cache interface {
	Fetch(key string, target interface{}, build func() (interface{}, error)) error
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

func (a *Analyze) Analyze() error {
	var rawFeatures [][]float64
	var result []int
	var err error

	a.cfg.Log("gathering features...")
	err = a.cache.Fetch("features", &rawFeatures,
		func() (interface{}, error) {
			return a.data.Features(), nil
		})
	if err != nil {
		return err
	}
	if len(rawFeatures) == 0 {
		return nil
	}

	a.cfg.Log("normalizing features...")
	p := newPayload(rawFeatures).
		apply(a.stats.Normalize(rawFeatures))

	a.cfg.Log("computing kmeans...")
	err = a.cache.Fetch("kmeans", &result,
		func() (interface{}, error) {
			return a.kmeans(p.data())
		})
	if err != nil {
		return err
	}

	a.cfg.Log("sorting dataset...")
	a.data.Sort(result)

	return nil
}

func (a *Analyze) kmeans(features [][]float64) ([]int, error) {
	return kmeans.Kmeans(features, a.cfg.Size(),
		a.stats.Distance, a.cfg.MaxIterations())
}
