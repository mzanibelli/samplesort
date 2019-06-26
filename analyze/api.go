package analyze

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
