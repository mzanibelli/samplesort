package sample_test

import (
	"reflect"
	"samplesort/sample"
	"testing"
)

func TestSample(t *testing.T) {
	t.Run("it should flatten data",
		func(t *testing.T) {
			input := map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": 0.001,
				},
			}
			expected := map[string]float64{"foo.bar": 0.001}
			SUT := new(sample.Sample)
			actual := SUT.Flatten(input)
			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	t.Run("it should extract floats",
		func(t *testing.T) {
			input := map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": 0.001,
					"baz": "hello",
				},
			}
			expected := map[string]float64{"foo.bar": 0.001}
			SUT := new(sample.Sample)
			actual := SUT.Flatten(input)
			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	t.Run("it should cache data",
		func(t *testing.T) {
			input := map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": 0.001,
				},
			}
			expected := map[string]float64{"foo.bar": 0.001}
			SUT := new(sample.Sample)
			SUT.Flatten(input)
			actual := SUT.Flatten()
			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
}
