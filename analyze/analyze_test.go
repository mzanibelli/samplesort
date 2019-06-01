package analyze_test

import (
	"samplesort/analyze"
	"testing"
)

func TestAnalyze(t *testing.T) {
	t.Run("it should process data and sort it afterwards",
		func(t *testing.T) {
			col := new(mockDataset)
			eng := new(mockEngine)
			cac := new(mockCache)
			SUT := analyze.New(col, eng, cac, 2, 0)
			if err := SUT.Analyze(); err != nil {
				t.Error("should not fail")
			}
			expected := 3
			actual := col.flag
			if expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
}

type mockDataset struct {
	flag int
}

func (d *mockDataset) Features() [][]float64 {
	d.flag++
	return [][]float64{
		{1, 2, 3},
		{4, 5, 6},
	}
}

func (d *mockDataset) Sort(centers []int) { d.flag += len(centers) }

type mockEngine struct{}

func (mockEngine) Compute([][]float64)                        {}
func (mockEngine) Normalize(data [][]float64) [][]float64     { return data }
func (mockEngine) Distance(s1, s2 []float64) (float64, error) { return 0, nil }

type mockCache struct{ err error }

func (m *mockCache) Fetch(
	key string,
	target interface{},
	build func() ([]byte, error),
) error {
	if m.err == nil {
		build()
	}
	return m.err
}
