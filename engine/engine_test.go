package engine_test

import (
	"samplesort/engine"
	"testing"
)

func TestEngine(t *testing.T) {
	cases := []struct {
		key    string
		input  float64
		output float64
	}{
		{"foo", 0, 0.5},
		{"bar", 0.5, 0.5},
		{"baz", 7.5, 0.75},
	}
	SUT := makeSUT()
	for _, c := range cases {
		t.Run(c.key, func(t *testing.T) {
			expected := c.output
			actual := SUT.Normalize(c.key, c.input)
			if expected != actual {
				t.Errorf("expected: %.10f, actual: %.10f", expected, actual)
			}
		})
	}
}

func makeSUT() *engine.Engine {
	e := engine.New(0.01)
	e.Update("foo", -1)
	e.Update("foo", 1)
	e.Update("bar", 0)
	e.Update("bar", 1)
	e.Update("baz", 0)
	e.Update("baz", 10)
	return e
}
