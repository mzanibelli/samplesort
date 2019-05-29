package analyze_test

import (
	"samplesort/analyze"
	"testing"
)

func TestAnalyze(t *testing.T) {
	t.Run("it should process data and sort it afterwards",
		func(t *testing.T) {
			data := new(mockDataset)
			SUT := analyze.New(data, 2, 0)
			SUT.Analyze()
			expected := 2
			actual := data.flag
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
	return [][]float64{{1, 2, 3}, {4, 5, 6}}
}

func (d *mockDataset) Sort([]int) { d.flag++ }
