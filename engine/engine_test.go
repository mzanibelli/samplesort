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
			result[i][j] = rand.Float64()
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
	seed := 42
	t.Run("input and output should have the same length",
		func(t *testing.T) {
			checkLength := func(g generator) bool {
				data := getData(g, seed)
				SUT := engine.New(mockConfig{})
				actual := SUT.Normalize(data)
				if len(data) != len(actual) {
					t.Log("input:", len(data))
					t.Log("output:", len(actual))
					return false
				}
				return true
			}
			if err := quick.Check(checkLength, nil); err != nil {
				t.Error("input and output are different")
			}
		})
	t.Run("output values should always be between 0 and 1",
		func(t *testing.T) {
			checkDecreasing := func(g generator) bool {
				data := getData(g, seed)
				SUT := engine.New(mockConfig{})
				actual := SUT.Normalize(data)
				for i, row := range data {
					for j := range row {
						if actual[i][j] < 0 || actual[i][j] > 1 {
							t.Log("input:", data[i][j])
							t.Log("output:", actual[i][j])
							return false
						}
					}
				}
				return true
			}
			if err := quick.Check(checkDecreasing, nil); err != nil {
				t.Error("invalid output range")
			}
		})
}

func TestDistance(t *testing.T) {
	t.Skip("not ready yet...")
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
