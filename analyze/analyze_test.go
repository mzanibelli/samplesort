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
			SUT := analyze.New(col, eng, 2, 0)
			SUT.Analyze()
			expected := 2
			actual := col.flag
			if expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
}

func TestDistance(t *testing.T) {
	cases := []struct {
		name  string
		error []float64
		input struct {
			i []float64
			j []float64
		}
		output struct {
			res float64
			err error
		}
	}{
		{
			name:  "obvious",
			error: []float64{1, 1, 10, 1},
			input: struct {
				i []float64
				j []float64
			}{
				[]float64{27, 12.3, 42.356, -2},
				[]float64{3, 12.0, 38.85, -1.7},
			},
			output: struct {
				res float64
				err error
			}{
				1,
				nil,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			col := new(mockDataset)
			eng := new(mockEngine)
			SUT := analyze.New(col, eng, 2, 0)
			res, err := SUT.Distance(c.error)(c.input.i, c.input.j)
			if c.output.res != res {
				t.Errorf("distance mismatch: expected: %v, actual: %v", c.output.res, res)
			}
			if c.output.err != err {
				t.Errorf("error mismatch: expected: %v, actual: %v", c.output.err, err)
			}
		})
	}
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

func (d *mockDataset) Sort([]int) { d.flag++ }

type mockEngine struct{}

func (mockEngine) Compute([][]float64) {}
func (mockEngine) SDs() []float64      { return []float64{} }
