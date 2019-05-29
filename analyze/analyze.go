package analyze

import (
	"github.com/bugra/kmeans"
)

type Analyze struct {
	data      Dataset
	dist      kmeans.DistanceFunction
	size      int
	threshold int
}

func (a *Analyze) DistanceFunction(f kmeans.DistanceFunction) { a.dist = f }

type Dataset interface {
	Features() [][]float64
	Sort([]int)
}

func New(data Dataset, size, threshold int) *Analyze {
	return &Analyze{
		data,
		kmeans.HammingDistance,
		size,
		threshold,
	}
}

func (a *Analyze) Analyze() {
	res, err := kmeans.Kmeans(a.data.Features(), a.size, a.dist, a.threshold)
	if err != nil {
		panic(err)
	}
	a.data.Sort(res)
}
