package analyze

import (
	"github.com/bugra/kmeans"
	"gonum.org/v1/gonum/mat"
)

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
	p := newPayload(rawFeatures)
	p.apply(a.stats.Normalize(rawFeatures))

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

// TODO: kmeans is not a deterministic algorithm so it shouldn't be
// unit-tested.
func (a *Analyze) kmeans(features [][]float64) ([]int, error) {
	return kmeans.Kmeans(features, a.cfg.Size(),
		a.stats.Distance, a.cfg.MaxIterations())
}

type payload struct {
	rows  int
	cols  int
	dense *mat.Dense
}

func newPayload(data [][]float64) *payload {
	if len(data) == 0 {
		return nil
	}
	r, c := len(data), len(data[0])
	p := &payload{r, c, mat.NewDense(r, c, make([]float64, r*c, r*c))}
	for i := range data {
		p.dense.SetRow(i, data[i])
	}
	return p
}

func (p *payload) data() [][]float64 {
	res := make([][]float64, p.rows, p.rows)
	for i := range res {
		row := make([]float64, p.cols, p.cols)
		copy(row, p.dense.RawRowView(i))
		res[i] = row
	}
	return res
}

func (p *payload) apply(f func(i, j int, v float64) float64) {
	p.dense.Apply(f, p.dense)
}

func (p *payload) r() int { return p.rows }
func (p *payload) c() int { return p.cols }

func (p *payload) String() string {
	return fmt.Sprintf("%v\n",
		mat.Formatted(p.dense, mat.Prefix(""), mat.Squeeze()))
}
