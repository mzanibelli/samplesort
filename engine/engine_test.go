package engine_test

import (
	"math/rand"
	"reflect"
	"samplesort/engine"
	"testing"
	"testing/quick"
)

var seed int64 = 42

type generator [][]float64

// Generate makes sure generated values are always slices of same-length
// slices of float64.
func (g generator) Generate(r *rand.Rand, size int) reflect.Value {
	g = make([][]float64, size, size)
	for i := range g {
		g[i] = make([]float64, size, size)
		for j := range g[i] {
			g[i][j] = rand.Float64() * float64(rand.Int()*10)
		}
	}
	return reflect.ValueOf(g)
}

func getData(g generator, seed int64) [][]float64 {
	input, ok := quick.Value(reflect.TypeOf(g),
		rand.New(rand.NewSource(seed)))
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
	t.Run("normalized values should be...?",
		func(t *testing.T) {
			t.Skip("normalization still not decided")
			checkSign := func(g generator) bool {
				data := getData(g, seed)
				engine.New().Normalize(data)
				return true
			}
			if err := quick.Check(checkSign, nil); err != nil {
				t.Error(err)
			}
		})
}

func TestDistance(t *testing.T) {
	if testing.Short() {
		return
	}
	t.Run("distance between same slices should be 0",
		func(t *testing.T) {
			checkRange := func(g generator) bool {
				data := getData(g, seed)
				SUT := engine.New()
				SUT.Normalize(data)
				for i := range data {
					input := data[i]
					if v, _ := SUT.Distance(input, input); v != 0 {
						t.Log("input:", data[i])
						return false
					}
				}
				return true
			}
			if err := quick.Check(checkRange, nil); err != nil {
				t.Error(err)
			}
		})
}
