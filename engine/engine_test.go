package engine_test

import (
	"math/rand"
	"reflect"
	"samplesort/engine"
	"testing"
	"testing/quick"
)

type generator [][]float64

// Generate makes sure generated values are always slices of same-length
// slices of float64.
func (result generator) Generate(r *rand.Rand, size int) reflect.Value {
	result = make([][]float64, size, size)
	for i := range result {
		result[i] = make([]float64, size, size)
		for j := range result[i] {
			result[i][j] = rand.Float64() * float64(rand.Int()*10)
		}
	}
	return reflect.ValueOf(result)
}

func getData(g generator, seed int) [][]float64 {
	input, ok := quick.Value(reflect.TypeOf(g),
		rand.New(rand.NewSource(42)))
	if !ok {
		panic("generator failed")
	}
	data, ok := input.Interface().(generator)
	if !ok {
		panic("type assertion failed")
	}
	return [][]float64(data)
}

func TestNormalize(t *testing.T) {
	if testing.Short() {
		return
	}
	seed := 42
	t.Run("output values should always be between 0 and 1",
		func(t *testing.T) {
			checkRange := func(g generator) bool {
				data := getData(g, seed)
				SUT := engine.New(mockConfig{}).Normalize(data)
				for i := range data {
					for j := range data[i] {
						input := data[i][j]
						output := SUT(i, j, input)
						if output < 0 || output > 1 {
							t.Log("input:", data[i][j])
							t.Log("output:", output)
							return false
						}
					}
				}
				return true
			}
			if err := quick.Check(checkRange, nil); err != nil {
				t.Error("invalid output range")
			}
		})
}

func TestDistance(t *testing.T) {
	t.Skip("TODO: find the correct way to compute distance")
	cases := []struct {
		name   string
		input  struct{ i, j []float64 }
		output float64
	}{
		{
			name: "obvious",
			input: struct{ i, j []float64 }{
				[]float64{5, 19.658, 42.356, -1256},
				[]float64{3, 12.0, 38.85, -1.7},
			},
			output: 2,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			SUT := engine.New(mockConfig{})
			SUT.Normalize([][]float64{
				{10, 20, 30, 40},
				{1, 2, 3, 4},
				// stds { 6.36, 12.72, 19.09, 25.45 }
			})
			expected := c.output
			actual, _ := SUT.Distance(c.input.i, c.input.j)
			if expected != actual {
				t.Errorf("distance mismatch: expected: %v, actual: %v", expected, actual)
			}
		})
	}
}

type mockConfig struct{}

func (mockConfig) MaxZScore() float64 { return 0.5 }
