package analyze

import "encoding/json"

var (
	TestOptionFeatures   [][]float64
	TestOptionCenters    []int
	TestOptionFetchError error
)

func MakeSUT() *Analyze {
	t := mock{}
	return New(t, t, t, t)
}

type mock struct{}

func (mock) Fetch(key string, target interface{}, build func() (interface{}, error)) error {
	switch target.(type) {
	case *[][]float64:
		swapJSON(TestOptionFeatures, target)
	case *[]int:
		swapJSON(TestOptionCenters, target)
	}
	return TestOptionFetchError
}

func (mock) Distance(s1, s2 []float64) (float64, error)              { return 0, nil }
func (mock) Features() [][]float64                                   { return nil }
func (mock) Log(vs ...interface{})                                   { return }
func (mock) MaxIterations() int                                      { return 10 }
func (mock) Normalize([][]float64) func(i, j int, v float64) float64 { return nop }
func (mock) Size() int                                               { return 5 }
func (mock) Sort([]int)                                              { return }

func nop(i, j int, v float64) float64 { return v }

func swapJSON(src, dst interface{}) {
	tmp, err := json.Marshal(src)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(tmp, dst)
	if err != nil {
		panic(err)
	}
}
