package engine_test

import (
	"math"
	"reflect"
	"samplesort/engine"
	"testing"
)

func TestEngine(t *testing.T) {
	cases := []struct {
		name   string
		input  [][]float64
		output []float64
	}{
		{
			name: "foo",
			input: [][]float64{
				{10, 1},
				{20, 2},
				{30, 3},
				{40, 4},
				{50, 5},
			},
			output: []float64{math.Sqrt(200), math.Sqrt(2)},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			SUT := makeSUT(c.input)
			expected := c.output
			actual := SUT.SDs()
			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	}
}

func makeSUT(data [][]float64) *engine.Engine {
	SUT := engine.New()
	SUT.Compute(data)
	return SUT
}
