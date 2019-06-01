package engine_test

import (
	"reflect"
	"samplesort/engine"
	"testing"
)

func TestNormalize(t *testing.T) {
	dataset := [][]float64{
		{1.3},
		{2.7},
		{0.25},
		{558.96},
		{1.3},
		{33.5},
		{0.25},
		{2.7},
		{-5863.22},
	}
	expected := [][]float64{
		{1.3},
		{2.7},
		{0.25},
		{388.63412492914335}, // should be much smaller
		{1.3},
		{33.5},
		{0.25},
		{2.7},
		{-883.2407522504036}, // and reduction should be even more effective here
	}
	SUT := engine.New()
	actual := SUT.Normalize(dataset)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestDistance(t *testing.T) {
	cases := []struct {
		name  string
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
			name: "obvious",
			input: struct {
				i []float64
				j []float64
			}{
				[]float64{5, 19.658, 42.356, -1256},
				[]float64{3, 12.0, 38.85, -1.7},
			},
			output: struct {
				res float64
				err error
			}{
				2,
				nil,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			SUT := engine.New()
			SUT.Normalize([][]float64{
				{10, 20, 30, 40},
				{1, 2, 3, 4},
				// stds { 6.36, 12.72, 19.09, 25.45 }
			})
			res, err := SUT.Distance(c.input.i, c.input.j)
			if c.output.res != res {
				t.Errorf("distance mismatch: expected: %v, actual: %v", c.output.res, res)
			}
			if c.output.err != err {
				t.Errorf("error mismatch: expected: %v, actual: %v", c.output.err, err)
			}
		})
	}
}
