package analyze_test

import (
	"encoding/json"
	"samplesort/analyze"
	"testing"
)

func TestAnalyze(t *testing.T) {
	t.Run("it should process data and sort it afterwards",
		func(t *testing.T) {
			col := new(mockDataset)
			eng := new(mockEngine)
			cac := new(mockCache)
			cfg := new(mockConfig)
			SUT := analyze.New(col, eng, cac, cfg)
			if err := SUT.Analyze(); err != nil {
				t.Error("should not fail")
			}
			expected := 3
			actual := col.flag
			if expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	t.Run("it should not fail if no data is found",
		func(t *testing.T) {
			col := &mockDataset{0, true}
			eng := new(mockEngine)
			cac := new(mockCache)
			cfg := new(mockConfig)
			SUT := analyze.New(col, eng, cac, cfg)
			if err := SUT.Analyze(); err != nil {
				t.Error("should not fail")
			}
		})
}

type mockDataset struct {
	flag    int
	nothing bool
}

var features = [][]float64{
	{1, 2, 3},
	{4, 5, 6},
}

var centers = []int{1, 2}

func (d *mockDataset) Features() [][]float64 {
	if d.nothing {
		return [][]float64{}
	}
	d.flag++
	return features
}

func (d *mockDataset) Sort(centers []int) {
	d.flag += len(centers)
}

type mockEngine struct{}

func (mockEngine) Normalize(data [][]float64) func(i, j int, v float64) float64 {
	return func(i, j int, v float64) float64 {
		return v
	}
}

func (mockEngine) Distance(s1, s2 []float64) (float64, error) { return 0, nil }

type mockCache struct{ err error }

func (m *mockCache) Fetch(
	key string,
	target interface{},
	build func() (interface{}, error),
) error {
	if m.err != nil {
		return m.err
	}
	build()
	switch target.(type) {
	case *[][]float64:
		unmarshal(features, target)
	case *[]int:
		unmarshal(centers, target)
	}
	return nil
}

func unmarshal(data, target interface{}) {
	tmp, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(tmp, target)
	if err != nil {
		panic(err)
	}
}

type mockConfig struct{}

func (mockConfig) Size() int             { return 5 }
func (mockConfig) MaxIterations() int    { return 10 }
func (mockConfig) Log(vs ...interface{}) {}
