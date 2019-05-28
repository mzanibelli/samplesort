package analyze

import (
	"github.com/bugra/kmeans"
)

const (
	size      int = 100
	threshold int = 0
)

type Analyze struct {
	done chan struct{}
	data Dataset
	dist kmeans.DistanceFunction
}

func (a *Analyze) DistanceFunction(f kmeans.DistanceFunction) { a.dist = f }

type Dataset interface {
	Features() [][]float64
	Sort([]int)
}

func New(data Dataset) *Analyze {
	return &Analyze{
		make(chan struct{}),
		data,
		kmeans.HammingDistance,
	}
}

func (a *Analyze) Analyze() <-chan struct{} {
	go func() {
		defer close(a.done)
		res, err := kmeans.Kmeans(
			a.data.Features(), size, a.dist, threshold)
		if err != nil {
			panic(err)
		}
		a.data.Sort(res)
	}()
	return a.done
}
